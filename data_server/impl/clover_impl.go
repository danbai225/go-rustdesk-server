package impl

import (
	"errors"
	logs "github.com/danbai225/go-logs"
	"github.com/ostafen/clover"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model"
)

const (
	objectIdField = "_id"
	TableNamePeer = "Peer"
)

type CloverDataSever struct {
	DB *clover.DB
}

func (c *CloverDataSever) InitDB() error {
	db, err := clover.Open("clover-db")
	if err != nil {
		logs.Err(err)
		return err
	}
	c.DB = db
	_ = c.DB.CreateCollection(TableNamePeer)
	return nil
}
func (c *CloverDataSever) GetPeerByUUID(uuid string) (*model.Peer, error) {
	first, err := c.DB.Query(TableNamePeer).Where(clover.Field("uuid").Eq(uuid)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	peer := model.Peer{}
	err = first.Unmarshal(&peer)
	return &peer, err
}

func (c *CloverDataSever) GetPeerByID(id string) (*model.Peer, error) {
	first, err := c.DB.Query(TableNamePeer).Where(clover.Field("id").Eq(id)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	peer := model.Peer{}
	err = first.Unmarshal(&peer)
	return &peer, err
}
func (c *CloverDataSever) AddPeer(peer *model.Peer) error {
	if peer == nil {
		return errors.New("nil peer")
	}
	p, err2 := c.GetPeerByUUID(peer.UUID)
	if err2 != nil {
		return err2
	} else if p != nil {
		return errors.New("exist peer")
	}
	m, err := common.ToMap(peer, "json")
	if err != nil {
		return err
	}
	document := clover.NewDocumentOf(m)
	_, err = c.DB.InsertOne(TableNamePeer, document)
	return err
}

func (c *CloverDataSever) GetPeerAll() ([]*model.Peer, error) {
	all, err := c.DB.Query(TableNamePeer).FindAll()
	if err != nil {
		return nil, err
	}
	peers := make([]*model.Peer, 0)
	for _, document := range all {
		p := &model.Peer{}
		_ = document.Unmarshal(p)
		peers = append(peers, p)
	}
	return peers, err
}
