package data_server

import (
	"go-rustdesk-server/data_server/impl"
	"go-rustdesk-server/model"
)

type DataSever interface {
	InitDB() error
	GetPeerByUUID(uuid string) (*model.Peer, error)
	GetPeerByID(id string) (*model.Peer, error)
	GetPeerAll() ([]*model.Peer, error)
	AddPeer(peer *model.Peer) error
	AddPeerOrUpdate(peer *model.Peer) error
	UpdatePeer(peer *model.Peer) error
	GetRelayByName(name string) (*model.Relay, error)
	UpdateRelay(relay *model.Relay) error
	AddRelay(relay *model.Relay) error
	AddRelayOrUpdate(relay *model.Relay) error
	GetRelayAllOnline() ([]*model.Relay, error)
	GetUserByName(name string) (*model.User, error)
	CheckToken(token string) (*model.User, error)
	GenToken(name string) (string, error)
	Close() error
}

var dataSever DataSever

func initDataSever() error {
	dataSever = new(impl.CloverDataSever)
	return dataSever.InitDB()
}
func GetDataSever() (DataSever, error) {
	if dataSever == nil {
		//new data SEVER
		err := initDataSever()
		return dataSever, err
	}
	return dataSever, nil
}
