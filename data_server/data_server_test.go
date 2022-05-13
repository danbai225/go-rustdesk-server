package data_server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model"
	"os"
	"testing"
)

func TestDB(t *testing.T) {
	sever, err := GetDataSever()
	if err != nil {
		t.Error(err)
	}
	logs.Info(sever)
}
func TestPeer(t *testing.T) {
	_ = os.RemoveAll("clover-db")
	sever, err := GetDataSever()
	if err != nil {
		t.Error(err)
	}
	err = sever.AddPeer(&model.Peer{
		ID:   "1",
		UUID: "3",
		PK:   []byte{1, 2, 3},
	})
	if err != nil {
		t.Error(err)
	}
	peer, err := sever.GetPeerByUUID("3")
	if err != nil {
		t.Error(err)
	}
	logs.Info(peer)
	peer.PK = []byte{3, 2, 1}
	err = sever.UpdatePeer(peer)
	if err != nil {
		t.Error(err)
	}
	all, err := sever.GetPeerAll()
	if err != nil {
		t.Error(err)
	}
	logs.Info(all)
}
