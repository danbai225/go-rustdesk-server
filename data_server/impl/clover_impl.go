package impl

import (
	"context"
	"errors"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/google/uuid"
	"github.com/ostafen/clover"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model"
	"strings"
	"sync"
	"time"
)

const (
	TableNamePeer  = "Peer"
	TableNameRelay = "Relay"
	TableNameUser  = "User"
	TableNameToken = "Token"
)

var cache = gcache.New()

type CloverDataSever struct {
	DB        *clover.DB
	peerLock  sync.RWMutex
	relayLock sync.RWMutex
	userLock  sync.RWMutex
}

func (c *CloverDataSever) Close() error {
	defer func() {
		c.peerLock.Unlock()
		c.relayLock.Unlock()
	}()
	c.peerLock.Lock()
	c.relayLock.Lock()
	return c.DB.Close()
}
func (c *CloverDataSever) InitDB() error {
	defer func() {
		c.peerLock.Unlock()
		c.relayLock.Unlock()
	}()
	c.peerLock.Lock()
	c.relayLock.Lock()
	db, err := clover.Open("clover-db")
	if err != nil {
		logs.Err(err)
		return err
	}
	c.DB = db
	_ = c.DB.CreateCollection(TableNamePeer)
	_ = c.DB.CreateCollection(TableNameRelay)
	_ = c.DB.CreateCollection(TableNameUser)
	_ = c.DB.CreateCollection(TableNameToken)
	_ = c.AddUser(&model.User{
		Name:     "admin",
		Password: "admin",
		IsAdmin:  true,
	})
	all, _ := c.getTokenAll()
	_ = cache.SetMap(context.Background(), all, time.Hour)
	return nil
}

func (c *CloverDataSever) getTokenAll() (map[interface{}]interface{}, error) {
	t, err := c.DB.Query(TableNameToken).FindFirst()
	if err != nil || t == nil {
		return nil, err
	}
	ts := make(map[interface{}]interface{})
	err = t.Unmarshal(&ts)
	return ts, err
}
func (c *CloverDataSever) saveToken(data map[interface{}]interface{}) error {
	if data == nil {
		return nil
	}
	_ = c.DB.Query(TableNameToken).Delete()
	m := make(map[string]interface{})
	for k, v := range data {
		m[k.(string)] = v
	}
	_, err := c.DB.InsertOne(TableNameToken, clover.NewDocumentOf(m))
	return err
}
func (c *CloverDataSever) GetPeerByUUID(uuid string) (*model.Peer, error) {
	defer c.peerLock.RUnlock()
	c.peerLock.RLock()
	first, err := c.DB.Query(TableNamePeer).Where(clover.Field("uuid").Eq(uuid)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	peer := model.Peer{}
	err = first.Unmarshal(&peer)
	return &peer, err
}

func (c *CloverDataSever) GetPeerByID(id string) (*model.Peer, error) {
	defer c.peerLock.RUnlock()
	c.peerLock.RLock()
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
	defer c.peerLock.Unlock()
	c.peerLock.Lock()
	m, err := common.ToMap(peer, "json")
	if err != nil {
		return err
	}
	document := clover.NewDocumentOf(m)
	_, err = c.DB.InsertOne(TableNamePeer, document)
	return err
}

func (c *CloverDataSever) GetPeerAll() ([]*model.Peer, error) {
	defer c.peerLock.RUnlock()
	c.peerLock.RLock()
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
	defer c.peerLock.Unlock()
	c.peerLock.Lock()
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
	defer c.relayLock.Unlock()
	c.relayLock.Lock()
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
	defer c.relayLock.Unlock()
	c.relayLock.Lock()
	m, err := common.ToMap(relay, "json")
	if err != nil {
		return err
	}
	return c.DB.Save(TableNameRelay, clover.NewDocumentOf(m))
}
func (c *CloverDataSever) GetRelayAllOnline() ([]*model.Relay, error) {
	c.relayLock.RLock()
	all, err := c.DB.Query(TableNameRelay).Where(clover.Field("online").Eq(true)).FindAll()
	if err != nil {
		return nil, err
	}
	c.relayLock.RUnlock()
	peers := make([]*model.Relay, 0)
	for _, document := range all {
		p := &model.Relay{}
		_ = document.Unmarshal(p)
		if p.LastRegTime.Add(time.Second * 60).After(time.Now()) {
			peers = append(peers, p)
		} else {
			p.Online = false
			err = c.UpdateRelay(p)
			if err != nil {
				return nil, err
			}
		}
	}
	return peers, err
}

func (c *CloverDataSever) GetRelayByName(name string) (*model.Relay, error) {
	defer c.relayLock.RUnlock()
	c.relayLock.RLock()
	first, err := c.DB.Query(TableNameRelay).Where(clover.Field("name").Eq(name)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	peer := model.Relay{}
	err = first.Unmarshal(&peer)
	return &peer, err
}

func (c *CloverDataSever) GetUserByName(name string) (*model.User, error) {
	defer c.userLock.RUnlock()
	c.userLock.RLock()
	first, err := c.DB.Query(TableNameUser).Where(clover.Field("name").Eq(name)).FindFirst()
	if err != nil || first == nil {
		return nil, err
	}
	user := model.User{}
	err = first.Unmarshal(&user)
	return &user, err
}
func (c *CloverDataSever) UpdateUser(user *model.User) error {
	defer c.userLock.Unlock()
	c.userLock.Lock()
	m, err := common.ToMap(user, "json")
	if err != nil {
		return err
	}
	return c.DB.Save(TableNameUser, clover.NewDocumentOf(m))
}
func (c *CloverDataSever) DelUser(name string) error {
	defer c.userLock.Unlock()
	c.userLock.Lock()
	err := c.DB.Query(TableNameUser).Where(clover.Field("name").Eq(name)).Delete()
	return err
}
func (c *CloverDataSever) GetUserAll() ([]*model.User, error) {
	defer c.userLock.RUnlock()
	c.userLock.RLock()
	all, err := c.DB.Query(TableNameUser).FindAll()
	if err != nil {
		return nil, err
	}
	users := make([]*model.User, 0)
	for _, document := range all {
		p := &model.User{}
		_ = document.Unmarshal(p)
		users = append(users, p)
	}
	return users, err
}
func (c *CloverDataSever) AddUser(user *model.User) error {
	if user == nil {
		return errors.New("nil user")
	}
	u, err := c.GetUserByName(user.Name)
	if err != nil {
		return err
	} else if u != nil {
		return errors.New("exist relay")
	}
	defer c.userLock.Unlock()
	c.userLock.Lock()
	m, err := common.ToMap(user, "json")
	if err != nil {
		return err
	}
	document := clover.NewDocumentOf(m)
	_, err = c.DB.InsertOne(TableNameUser, document)
	return err
}
func (c *CloverDataSever) CheckToken(token string) (*model.User, error) {
	defer c.userLock.RUnlock()
	c.userLock.RLock()
	val, err := cache.Get(context.Background(), token)
	if err != nil || val == nil {
		return nil, errors.New("auth error")
	}
	user, err := c.GetUserByName(val.String())
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (c *CloverDataSever) GenToken(name string) (string, error) {
	defer c.userLock.Unlock()
	c.userLock.Unlock()
	user, err := c.GetUserByName(name)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not exist")
	}
	token := strings.ReplaceAll(uuid.New().String(), "-", "")
	err = cache.Set(context.Background(), token, name, time.Hour*256)
	if err != nil {
		return "", err
	}
	data, _ := cache.Data(context.Background())
	_ = c.saveToken(data)
	return token, nil
}
func (c *CloverDataSever) DelToken(token string) error {
	defer c.userLock.Unlock()
	c.userLock.Lock()
	_, err := cache.Remove(context.Background(), token)
	if err != nil {
		return err
	}
	data, _ := cache.Data(context.Background())
	_ = c.saveToken(data)
	return nil
}
