package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
)

func RendezvousMessageRegisterPeer(message *model_proto.RegisterPeer) *model_proto.RegisterPeerResponse {
	res := &model_proto.RegisterPeerResponse{}
	peer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if peer.UUID != "" {
		return res
	} else {
		res.RequestPk = true
	}
	return res
}
func RendezvousMessageRegisterPk(message *model_proto.RegisterPk) *model_proto.RegisterPkResponse {
	res := &model_proto.RegisterPkResponse{Result: model_proto.RegisterPkResponse_SERVER_ERROR}
	peer, err := dataSever.GetPeerByID(message.GetId())
	if err != nil {
		logs.Err(err)
		return res
	}
	if peer.UUID != "" {
		if peer.UUID != string(message.GetUuid()) {
			res.Result = model_proto.RegisterPkResponse_UUID_MISMATCH
			return res
		}
	}
	return res
}
