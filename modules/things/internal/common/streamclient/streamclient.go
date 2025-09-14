package streamclient

import (
	"time"

	"github.com/lostvip-com/lv_framework/lv_log"
)

const (
	pubTimeout = time.Millisecond * 10
)

type streamClient struct {
	msgCh chan RpcData
}

func (c *streamClient) Send(data RpcData) {
	select {
	case c.msgCh <- data:
	case <-time.After(pubTimeout):
		lv_log.Warnf("send stream message timeout, data: %+v", data)
	}
}

func (c *streamClient) Recv() <-chan RpcData {
	return c.msgCh
}

func NewStreamClient() *streamClient {
	return &streamClient{
		msgCh: make(chan RpcData),
	}
}
