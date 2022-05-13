package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func handlerMsg(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
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
		response = model_proto.NewRendezvousMessage(RendezvousMessagePunchHoleRequest(HoleRequest))
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
			return
		}
		_, err2 = writer.Write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	}
}
func handlerMsgUDP(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil {
		logs.Err(err)
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRegisterPk:
		//注册公钥
		RegisterPk := message.GetRegisterPk()
		if RegisterPk == nil {
			return
		}
		pk := RendezvousMessageRegisterPk(RegisterPk)
		response = model_proto.NewRendezvousMessage(pk)
		if pk.Result == model_proto.RegisterPkResponse_OK {
			if _, ok := connPeerMap[RegisterPk.GetId()]; !ok {
				connPeerMap[RegisterPk.GetId()] = writer
			}
		}

	case model_proto.TypeRendezvousMessageRegisterPeer:
		//注册id
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
		response = model_proto.NewRendezvousMessage(RendezvousMessageSoftwareUpdate(SoftwareUpdate))
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
			return
		}
		_, err2 = writer.Write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	}
}
