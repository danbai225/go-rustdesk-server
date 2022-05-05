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
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRegisterPk:
		RegisterPk := message.GetRegisterPk()
		if RegisterPk == nil {
			return
		}
		response := model_proto.NewRendezvousMessage(&model_proto.RendezvousMessage_RegisterPkResponse{
			RegisterPkResponse: &model_proto.RegisterPkResponse{
				Result: 0,
			},
		})
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
		}
		err2 = write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
}
