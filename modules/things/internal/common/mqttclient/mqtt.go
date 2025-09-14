package mqttclient

import (
	"context"
	"things/internal/common/dtos"
)

type MQTTClient interface {
	RegisterConnectCallback(dtos.ConnectHandler)
	RegisterDisconnectCallback(dtos.CallbackHandler)
	AsyncPublish(ctx context.Context, topic string, payload []byte, isSync bool)
	Close()
	GetConnectStatus() bool
}
