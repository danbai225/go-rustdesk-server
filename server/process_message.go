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

func RendezvousMessageRegisterPeer(message *model_proto.RegisterPeer, writer *common.Writer) *model_proto.RegisterPeerResponse {
	res := &model_proto.RegisterPeerResponse{}
	peer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if peer == nil {
		res.RequestPk = true
	} else {
		ipChange := false
		res.RequestPk = false
		w, err1 := common.GetWriter(message.GetId(), common.UDP)
		if err1 != nil {
			ipChange = true
		} else if w.GetAddrStr() != writer.GetAddrStr() {
			ipChange = true
		}
		now := time.Now()
		res.RequestPk = len(peer.PK) == 0 || ipChange
		if ipChange {
			peer.IP = writer.GetAddr().GetIP()
			peer.LastRegTime = &now
			err = dataSever.UpdatePeer(peer)
			if err != nil {
				logs.Err(err)
				return res
			}
			writer.SetKey(message.GetId())
		}
	}
	return res
}
func RendezvousMessageRegisterPk(message *model_proto.RegisterPk, writer *common.Writer) *model_proto.RegisterPkResponse {
	res := &model_proto.RegisterPkResponse{Result: model_proto.RegisterPkResponse_SERVER_ERROR}
	if len(message.GetId()) < common.MinKeyLen {
		res.Result = model_proto.RegisterPkResponse_UUID_MISMATCH
		return res
	}
	//改变 ID
	if id := message.GetOldId(); id != "" {
		idPeer, err := dataSever.GetPeerByID(id)
		if err != nil {
			logs.Debug(err)
			return res
		}
		idPeer.ID = message.Id
		err = dataSever.UpdatePeer(idPeer)
		if err != nil {
			res.Result = model_proto.RegisterPkResponse_SERVER_ERROR
			logs.Err(err)
		} else {
			res.Result = model_proto.RegisterPkResponse_OK
			writer.SetKey(message.Id)
		}
		return res
	}

	idPeer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Debug(err)
		return res
	}
	change := false
	if idPeer != nil {
		if idPeer.UUID == "" {
			change = true
		} else if idPeer.UUID == string(message.GetUuid()) {
			//存在注册
			if string(idPeer.PK) != string(message.GetPk()) {
				if idPeer.IP != writer.GetAddr().GetIP() {
					//都不匹配
					res.Result = model_proto.RegisterPkResponse_UUID_MISMATCH
					return res
				}
				change = true
			}
		} else {
			res.Result = model_proto.RegisterPkResponse_ID_EXISTS
			return res
		}
		res.Result = model_proto.RegisterPkResponse_OK
	}
	getWriter, err := common.GetWriter(message.GetId(), common.UDP)
	ipChange := false
	if err == nil {
		if getWriter.GetAddrStr() != writer.GetAddrStr() {
			ipChange = true
		}
	}
	if ipChange {
		idPeer.IP = writer.GetAddr().GetIP()
		writer.SetKey(idPeer.ID)
	}
	change = ipChange || change || idPeer == nil
	if change {
		now := time.Now()
		peer := &model.Peer{
			ID:          message.Id,
			UUID:        string(message.Uuid),
			PK:          message.Pk,
			LastRegTime: &now,
		}
		err = dataSever.AddPeerOrUpdate(peer)
		if err != nil {
			res.Result = model_proto.RegisterPkResponse_SERVER_ERROR
			logs.Err(err)
		} else {
			res.Result = model_proto.RegisterPkResponse_OK
		}
	}
	if res.Result == model_proto.RegisterPkResponse_OK {
		writer.SetKey(message.GetId())
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
		logs.Debug(err)
		res.OtherFailure = err.Error()
		return res
	}
	if peer == nil {
		res.Failure = model_proto.PunchHoleResponse_ID_NOT_EXIST
		return res
	}
	w, err := common.GetWriter(message.GetId(), common.UDP)
	if err != nil {
		res.Failure = model_proto.PunchHoleResponse_OFFLINE
	} else {
		rendezvousMessage := model_proto.NewRendezvousMessage(&model_proto.FetchLocalAddr{
			SocketAddr:  my_bytes.EncodeAddr(writer.GetAddrStr()),
			RelayServer: rendezvousServers[0],
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

			res.Pk = common.GetSignPK(m.GetVersion(), peer.ID, peer.PK)
			res.Union = &model_proto.PunchHoleResponse_IsLocal{IsLocal: true}
		}
	}
	return res
}
func RendezvousMessageTestNatRequest(message *model_proto.TestNatRequest, writer *common.Writer) *model_proto.TestNatResponse {
	res := &model_proto.TestNatResponse{
		Port: int32(writer.GetAddr().GetPort()),
	}
	res.Cu = &model_proto.ConfigUpdate{
		Serial:            message.Serial,
		RendezvousServers: rendezvousServers,
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
	w, err := common.GetWriter(message.GetId(), common.UDP)
	if err != nil {
		return nil
	} else {
		peer, err1 := dataSever.GetPeerByID(message.Id)
		if err1 != nil {
			logs.Debug(err1)
			return res
		}
		w.SendMsg(model_proto.NewRendezvousMessage(message))
		lMsg := getMsgForm(message.GetId(), model_proto.TypeRendezvousMessageRelayResponse, 3)
		if lMsg == nil {
			return res
		} else if m, ok1 := lMsg.(*model_proto.RelayResponse); ok1 {
			m.Union = &model_proto.RelayResponse_Pk{
				Pk: common.GetSignPK(m.GetVersion(), m.GetId(), peer.PK),
			}
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
func ConfigureUpdate(writer *common.Writer) {
	writer.SendMsg(model_proto.NewRendezvousMessage(&model_proto.ConfigUpdate{
		Serial:            serial,
		RendezvousServers: rendezvousServers,
	}))
}

func RendezvousMessagePunchHoleSent(message *model_proto.PunchHoleSent, writer *common.Writer) *model_proto.PunchHoleResponse {
	peer, err := dataSever.GetPeerByID(message.Id)
	if err != nil {
		logs.Debug(err)
		return nil
	}
	res := &model_proto.PunchHoleResponse{
		SocketAddr:  my_bytes.EncodeAddr(writer.GetAddrStr()),
		Pk:          common.GetSignPK(message.GetVersion(), message.GetId(), peer.PK),
		RelayServer: rendezvousServers[0],
		Union: &model_proto.PunchHoleResponse_NatType{
			NatType: message.NatType,
		},
	}
	marshal, _ := proto.Marshal(res)
	addr := my_bytes.DecodeAddr(message.GetSocketAddr())
	_, err = writer.WriteToAddr(marshal, addr)
	if err != nil {
		logs.Err(err)
	}
	return res
}
func RendezvousMessageConfigureUpdate(message *model_proto.ConfigUpdate) {
	logs.Debug(message.Serial, message.RendezvousServers)
}
