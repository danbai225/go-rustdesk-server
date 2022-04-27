package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func Start() {
	go common.NewMonitor("udp", ":21116", handlerMsgUDP).Start()
	common.NewMonitor("tcp", ":21116", handlerMsg).Start()
}
func handlerMsg(messageData []byte, write func(data []byte) error) {
	message := model_proto.Message{}
	err := proto.Unmarshal(messageData, &message)
	if err != nil {
		logs.Err(err)
	}
	logs.Info(message.Union)
}
func handlerMsgUDP(messageData []byte, write func(data []byte) error) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(messageData, &message)
	if err != nil {
		logs.Err(err)
	}
	switch reflect.TypeOf(message.Union).Kind() {
	case model_proto.Type_RendezvousMessage_RegisterPk:
		RegisterPk := message.GetRegisterPk()
		logs.Info(string(RegisterPk.Uuid))
		response := &model_proto.RegisterPkResponse{
			Result: 0,
		}
		response.ProtoReflect()
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
		}
		err2 = write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}

	}
}
