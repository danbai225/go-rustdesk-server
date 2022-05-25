package impl

import (
	"errors"
	logs "github.com/danbai225/go-logs"
	"github.com/ostafen/clover"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model"
)

const (
	objectIdField  = "_id"
	TableNamePeer  = "Peer"
	TableNameRelay = "Relay"
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

func (c *CloverDataSever) UpdatePeer(peer *model.Peer) error {
	m, err := common.ToMap(peer, "json")
	if err != nil {
		return err
	}
	return c.DB.Save(TableNamePeer, clover.NewDocumentOf(m))
}
func (c *CloverDataSever) AddPeerOrUpdate(peer *model.Peer) error {
	if err := c.AddPeer(peer); err != nil {
		return c.UpdatePeer(peer)
	}
	return nil
}
func (c *CloverDataSever) DelPeerByUUID(uuid string) error {
	peer, err := c.GetPeerByUUID(uuid)
	if err == nil {
		return c.DB.Query(TableNamePeer).Where(clover.Field("_id").Eq(peer.Uid)).Delete()
	}
	return nil
}

func (c *CloverDataSever) AddRelay(relay *model.Relay) error {
	if relay == nil {
		return errors.New("nil relay")
	}
	p, err2 := c.GetRelayByName(relay.Name)
	if err2 != nil {
		return err2
	} else if p != nil {
		return errors.New("exist relay")
	}
	m, err := common.ToMap(relay, "json")
	if err != nil {
		return err
	}
	document := clover.NewDocumentOf(m)
	_, err = c.DB.InsertOne(TableNameRelay, document)
	return err
}
func (c *CloverDataSever) AddRelayOrUpdate(relay *model.Relay) error {
	if err := c.AddRelay(relay); err != nil {
		return c.UpdateRelay(relay)
	}
	return nil
}
func (c *CloverDataSever) UpdateRelay(relay *model.Relay) error {
	m, err := common.ToMap(relay, "json")
	if err != nil {
		return err
	}
	return c.DB.Save(TableNameRelay, clover.NewDocumentOf(m))
}
func (c *CloverDataSever) GetRelayAllOnline() ([]*model.Relay, error) {
	all, err := c.DB.Query(TableNameRelay).Where(clover.Field("online").Eq(true)).FindAll()
	if err != nil {
		return nil, err
	}
	peers := make([]*model.Relay, 0)
	for _, document := range all {
		p := &model.Relay{}
		_ = document.Unmarshal(p)
		peers = append(peers, p)
	}
	return peers, err
}

func (c *CloverDataSever) GetRelayByName(name string) (*model.Relay, error) {
	first, err := c.DB.Query(TableNameRelay).Where(clover.Field("name").Eq(name)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	peer := model.Relay{}
	err = first.Unmarshal(&peer)
	return &peer, err
}
