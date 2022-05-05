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
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeMessageMessageCliprdr:
		Cliprdr := message.GetCliprdr()
		if Cliprdr == nil {
			return
		}
		if Cliprdr.Union != nil {
			switch reflect.TypeOf(Cliprdr.Union).String() {
			default:
				logs.Info(reflect.TypeOf(Cliprdr.Union).String())
			}
		}
	//response = model_proto.NewRendezvousMessage(&model_proto.RegisterPkResponse{Result: model_proto.RegisterPkResponse_OK})
	case model_proto.TypeMessageMessageLoginResponse:
		LoginResponse := message.GetLoginResponse()
		if LoginResponse == nil {
			return
		}
		logs.Info(LoginResponse.GetError())
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
			return
		}
		err2 = write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	}
}
func handlerMsgUDP(messageData []byte, write func(data []byte) error) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(messageData, &message)
	if err != nil {
		logs.Err(err)
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRegisterPk:
		RegisterPk := message.GetRegisterPk()
		if RegisterPk == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(&model_proto.RegisterPkResponse{Result: model_proto.RegisterPkResponse_OK})
	case model_proto.TypeRendezvousMessageRegisterPeer:
		RegisterPeer := message.GetRegisterPeer()
		if RegisterPeer == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(&model_proto.RegisterPeerResponse{
			RequestPk: true,
		})
	case model_proto.TypeRendezvousMessageSoftwareUpdate:
		//软件更新
		SoftwareUpdate := message.GetSoftwareUpdate()
		if SoftwareUpdate == nil {
			return
		}
		logs.Info(SoftwareUpdate.Url)
		response = model_proto.NewRendezvousMessage(&model_proto.SoftwareUpdate{
			Url: "",
		})
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
			return
		}
		err2 = write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	}
}
