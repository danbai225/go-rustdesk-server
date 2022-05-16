package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model"
	"go-rustdesk-server/model/model_proto"
	"go-rustdesk-server/my_bytes"
	"google.golang.org/protobuf/proto"
	"time"
)

func RendezvousMessageRegisterPeer(message *model_proto.RegisterPeer) *model_proto.RegisterPeerResponse {
	res := &model_proto.RegisterPeerResponse{}
	peer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if peer == nil {
		res.RequestPk = true
	} else if len(peer.PK) == 0 {
		res.RequestPk = true
	}
	return res
}
func RendezvousMessageRegisterPk(message *model_proto.RegisterPk) *model_proto.RegisterPkResponse {
	res := &model_proto.RegisterPkResponse{Result: model_proto.RegisterPkResponse_SERVER_ERROR}
	idPeer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if idPeer != nil {
		if idPeer.UUID == string(message.GetUuid()) {
			res.Result = model_proto.RegisterPkResponse_OK
			return res
		}
		res.Result = model_proto.RegisterPkResponse_ID_EXISTS
		return res
	}
	res.Result = model_proto.RegisterPkResponse_OK
	peer := &model.Peer{
		ID:   message.Id,
		UUID: string(message.Uuid),
		PK:   message.Pk,
	}
	err = dataSever.AddPeer(peer)
	if err != nil {
		res.Result = model_proto.RegisterPkResponse_SERVER_ERROR
		logs.Err(err)
	} else {
		res.Result = model_proto.RegisterPkResponse_OK
	}
	return res
}
func RendezvousMessageSoftwareUpdate(message *model_proto.SoftwareUpdate) *model_proto.SoftwareUpdate {
	res := &model_proto.SoftwareUpdate{}
	return res
}
func RendezvousMessagePunchHoleRequest(message *model_proto.PunchHoleRequest, writer *common.Writer) *model_proto.PunchHoleResponse {
	res := &model_proto.PunchHoleResponse{}
	peer, err := dataSever.GetPeerByID(message.Id)
	if err != nil {
		logs.Err(err)
		res.OtherFailure = err.Error()
		return res
	}
	if peer == nil {
		res.Failure = model_proto.PunchHoleResponse_ID_NOT_EXIST
	}
	w, err := common.GetWriter(message.GetId(), "udp")
	if err != nil {
		res.Failure = model_proto.PunchHoleResponse_OFFLINE
	} else {
		rendezvousMessage := model_proto.NewRendezvousMessage(&model_proto.FetchLocalAddr{
			SocketAddr:  my_bytes.EncodeAddr(writer.GetAddrStr()),
			RelayServer: "1.14.47.89",
		})
		marshal, err2 := proto.Marshal(rendezvousMessage)
		if err2 != nil {
			logs.Err(err2)
			return res
		}
		_, err2 = w.Write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
		lMsg := getMsgForm(message.GetId(), model_proto.TypeRendezvousMessageLocalAddr, 3)
		if lMsg == nil {
			res.OtherFailure = "NoReturnMessage"
			return res
		}
		if m, ok1 := lMsg.(*model_proto.LocalAddr); ok1 {
			res.SocketAddr = m.GetLocalAddr()
			res.RelayServer = m.GetRelayServer()
			//pk := model_proto.IdPk{
			//	Id: peer.ID,
			//	Pk: peer.PK,
			//}
			//data, _ := proto.Marshal(&pk)
			res.Pk = peer.PK
			res.Union = &model_proto.PunchHoleResponse_IsLocal{IsLocal: true}
		}
	}
	return res
}
func RendezvousMessageTestNatRequest(message *model_proto.TestNatRequest) *model_proto.TestNatResponse {
	res := &model_proto.TestNatResponse{}
	res.Cu = &model_proto.ConfigUpdate{
		Serial:            message.Serial,
		RendezvousServers: []string{"192.168.0.110"},
	}
	return res
}
func RendezvousMessageLocalAddr(message *model_proto.LocalAddr, writer *common.Writer) {
	r.Put(&ringMsg{
		ID:      message.GetId(),
		Type:    model_proto.TypeRendezvousMessageLocalAddr,
		TimeOut: 3,
		InsTime: time.Now(),
		Val:     message,
	})
}
func RendezvousMessageRequestRelay(message *model_proto.RequestRelay) *model_proto.RelayResponse {
	res := &model_proto.RelayResponse{}
	w, err := common.GetWriter(message.GetId(), "udp")
	if err != nil {
		return nil
	} else {
		rendezvousMessage := model_proto.NewRendezvousMessage(message)
		marshal, err2 := proto.Marshal(rendezvousMessage)
		if err2 != nil {
			return nil
		}
		_, _ = w.Write(marshal)
		lMsg := getMsgForm(message.GetId(), model_proto.TypeRendezvousMessageRelayResponse, 3)
		if lMsg == nil {
			return res
		} else if m, ok1 := lMsg.(*model_proto.RelayResponse); ok1 {
			res = m
		}
	}
	return res
}
func RendezvousMessageRelayResponse(message *model_proto.RelayResponse) {
	r.Put(&ringMsg{
		ID:      message.GetId(),
		Type:    model_proto.TypeRendezvousMessageRelayResponse,
		TimeOut: 3,
		InsTime: time.Now(),
		Val:     message,
	})
}
