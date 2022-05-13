package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model"
	"go-rustdesk-server/model/model_proto"
)

func RendezvousMessageRegisterPeer(message *model_proto.RegisterPeer) *model_proto.RegisterPeerResponse {
	res := &model_proto.RegisterPeerResponse{}
	peer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if peer == nil {
		return res
	} else {
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
	uuPeer, err := dataSever.GetPeerByUUID(string(message.GetUuid()))
	if err != nil {
		logs.Err(err)
		return res
	}
	if idPeer != nil {
		res.Result = model_proto.RegisterPkResponse_ID_EXISTS
		return res
	}
	if uuPeer != nil {
		res.Result = model_proto.RegisterPkResponse_UUID_MISMATCH
		return res
	}
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
func RendezvousMessagePunchHoleRequest(message *model_proto.PunchHoleRequest) *model_proto.PunchHoleResponse {
	res := &model_proto.PunchHoleResponse{}
	return res
}
