package mqttclient

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"things/internal/common/dtos"
	"things/internal/common/errort"
	"things/internal/common/middleware"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/lostvip-com/lv_framework/lv_log"
	pkgerr "github.com/pkg/errors"
	"go.uber.org/atomic"
)

var (
	cache sync.Map
)

type message struct {
	ctx     context.Context
	isSync  bool
	topic   string
	payload []byte
}

type mqttClient struct {
	ctx               context.Context
	cancel            context.CancelFunc
	clientId          string
	client            mqtt.Client
	connectHandler    dtos.ConnectHandler
	disconnectHandler dtos.CallbackHandler
	messageChannel    chan message
	status            *atomic.Bool
	pubTopic          string
}

func NewMQTTClient(req dtos.NewMQTTClient,
	consumeCallback func(mqtt.Client, mqtt.Message),
	connF dtos.ConnectHandler,
	disConnF dtos.CallbackHandler) (MQTTClient, error) {

	mqttCli, exist := cache.Load(req.ClientId)
	if exist {
		return mqttCli.(MQTTClient), nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	var mc = &mqttClient{
		ctx:               ctx,
		cancel:            cancel,
		clientId:          req.ClientId,
		connectHandler:    connF,
		disconnectHandler: disConnF,
		status:            atomic.NewBool(false),
		messageChannel:    make(chan message, 200),
	}
	opts := mqtt.NewClientOptions().SetClientID(req.ClientId).SetUsername(req.Username).SetPassword(req.Password)
	opts = opts.SetAutoReconnect(true).SetCleanSession(true).SetKeepAlive(5 * time.Second).SetMaxReconnectInterval(10 * time.Second).
		SetConnectRetry(true).SetConnectRetryInterval(time.Second)
	w := os.Stdout
	opts = opts.SetTLSConfig(&tls.Config{
		KeyLogWriter: w,
		//InsecureSkipVerify: true,	//如果本地出现证书过期，可以打开这个选项
	})

	// add mqttclient internal log
	mqtt.ERROR = log.New(os.Stdout, "", log.LstdFlags)

	opts.AddBroker(req.Broker)
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		defer func() {
			if err := recover(); err != nil {
				err = pkgerr.WithStack(pkgerr.New(fmt.Sprintf("internel err: %v", err)))
				lv_log.Errorf("internal error: %s,%s", err, string(debug.Stack()))
			}
		}()

		for _, topic := range req.SubTopics {
			token := client.Subscribe(topic, 1, consumeCallback)
			token.Wait()
			if err := token.Error(); err != nil {
				lv_log.Errorf("failed to consume message: %s", err.Error())
				continue
			}
			lv_log.Infof("Subscribe success,topic: %s", topic)
		}

		// business process
		mc.setStatus(true)
		//ctx := middleware.WithCorrelationId(context.Background())  ?????????????
		//lv_log.Info("gateway online, mqttclient connected!", middleware.FromContext(ctx))

		// callback connectHandler
		if mc.connectHandler != nil {
			mc.connectHandler(ctx)
		}
	}).SetConnectionLostHandler(func(client mqtt.Client, err error) {
		mc.setStatus(false)
		//ctx := middleware.WithCorrelationId(context.Background())  ?????

		err = pkgerr.Wrap(err, "gateway offline, mqttclient disconnect")
		err = errort.NewCommonErr(errort.MqttConnFail, err)

		//lv_log.Errorf("%v, %v", err.Error(), middleware.FromContext(ctx))
		if mc.disconnectHandler != nil {
			mc.disconnectHandler(ctx, dtos.CallbackMessage{
				Error: err,
			})
		}
	})
	opts.SetDefaultPublishHandler(consumeCallback)

	var err error
	mc.client = mqtt.NewClient(opts)
	token := mc.client.Connect()

	if ok := token.WaitTimeout(time.Second * 10); !ok {
		err = errort.NewCommonErr(errort.MqttConnFail, pkgerr.New("mqttclient client connect timeout"))
	}

	if err1 := token.Error(); err1 != nil {
		err = fmt.Errorf("%v, %v", err1, err)
		err = errort.NewCommonErr(errort.MqttConnFail, err)
	}
	mc.pubTopic = req.PubTopic
	go mc.run()
	cache.Store(req.ClientId, mc)
	return mc, err
}

func (c *mqttClient) Close() {
	c.cancel()
	c.client.Disconnect(3000)
	cache.Delete(c.clientId)
}

func (c *mqttClient) RegisterConnectCallback(cb dtos.ConnectHandler) {
	if c.connectHandler == nil {
		c.connectHandler = cb
	}
}

func (c *mqttClient) RegisterDisconnectCallback(cb dtos.CallbackHandler) {
	if c.disconnectHandler == nil {
		c.disconnectHandler = cb
	}
}

func (c *mqttClient) setStatus(status bool) {
	c.status.Store(status)
}

func (c *mqttClient) getStatus() bool {
	return c.status.Load()
}

func (c *mqttClient) publish(msg message) (err error) {
	token := c.client.Publish(msg.topic, 1, false, msg.payload)
	if !msg.isSync {
		return
	}

	go func() {
		// error handing https://github.com/eclipse/paho.mqtt.golang#error-handling
		_ = token.WaitTimeout(3 * time.Second)
		if token.Error() != nil {
			lv_log.Errorf("failed to publish %s, err: %v, %v", msg.topic, token.Error(), middleware.FromContext(msg.ctx))
			return
		}
		lv_log.Infof("sync publish to %s message success, %v", msg.topic, middleware.FromContext(msg.ctx))
	}()
	return
}

func (c *mqttClient) run() {
	for {
		select {
		case <-c.ctx.Done():
			close(c.messageChannel)
			lv_log.Error("mqttclient client is closed")
			return
		case msg := <-c.messageChannel:
			if c.getStatus() {
				// 消息发送
				if err := c.publish(msg); err != nil {
					lv_log.Error("failed to publish message, topic: %v, %v", msg.topic, middleware.FromContext(msg.ctx))
				}
			} else {
				// 消息丢弃
				lv_log.Warnf("failed to publish, mqttclient client is reconnecting, message must be lost, topic: %v, %v", msg.topic, middleware.FromContext(msg.ctx))
				if len(c.messageChannel) < cap(c.messageChannel) {
					// 等待 100ms 缓冲
					time.After(100 * time.Millisecond)
				}
			}
		}
	}
}

// 异步发送，实际上是发送到channel中排队
func (c *mqttClient) AsyncPublish(ctx context.Context, topic string, payload []byte, isSync bool) {
	//c.messageChannel <- message{
	//	ctx:     ctx,
	//	isSync:  isSync,
	//	topic:   topic,
	//	payload: payload,
	//}
	if c.getStatus() {
		_ = c.publish(message{
			ctx:     ctx,
			isSync:  isSync,
			topic:   topic,
			payload: payload,
		})
	} else {
		lv_log.Warnf("failed to publish, mqttclient client is reconnecting, message must be lost, topic: %v, %v", topic, middleware.FromContext(ctx))
	}

}

func (c *mqttClient) GetConnectStatus() bool {
	return c.getStatus()
}
