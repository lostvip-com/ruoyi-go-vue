package streamclient

import (
	"encoding/json"
)

type StreamClient interface {
	Send(data RpcData)
	Recv() <-chan RpcData
}

type RpcData struct {
	Code    int32
	ReqId   string
	ErrCode uint32
	Data    interface{}
}

func (rd *RpcData) String() string {
	str, _ := json.Marshal(rd)
	return string(str)
}
