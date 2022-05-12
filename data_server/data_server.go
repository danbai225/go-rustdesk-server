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
