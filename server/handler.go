package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func handlerMsg(messageData []byte, write func(data []byte) error) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(messageData, &message)
	if err != nil {
		logs.Err(err)
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessagePunchHoleRequest:
		HoleRequest := message.GetPunchHoleRequest()
		if HoleRequest == nil {
			return
		}
		logs.Info(HoleRequest.Id, HoleRequest.NatType, HoleRequest.ConnType)
		response = model_proto.NewRendezvousMessage(&model_proto.PunchHoleResponse{
			SocketAddr:   nil,
			Pk:           nil,
			Failure:      0,
			RelayServer:  "",
			Union:        nil,
			OtherFailure: "",
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
		response = model_proto.NewRendezvousMessage(RendezvousMessageRegisterPeer(RegisterPeer))
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
