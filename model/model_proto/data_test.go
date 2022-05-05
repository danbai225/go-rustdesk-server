package model_proto

import (
	logs "github.com/danbai225/go-logs"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestOK(t *testing.T) {
	message := RendezvousMessage{}
	message.Union = &RendezvousMessage_RegisterPkResponse{
		RegisterPkResponse: &RegisterPkResponse{
			Result: 0,
		},
	}
	marshal, err := proto.Marshal(&message)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(marshal)
	message2 := RendezvousMessage{}
	err = proto.Unmarshal(marshal, &message2)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(message2.Union)
}
