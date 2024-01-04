package server

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/v2/container/gring"
	"go-rustdesk-server/api_server"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
	"go-rustdesk-server/model"
	"math/rand"
)

var dataSever data_server.DataSever
var r = gring.New(32, true)
var rendezvousServers = make([]string, 0)
var serial = int32(1)

func Start() {
	common.LoadKey()
	var err error
	dataSever, err = data_server.GetDataSever()
	if err != nil {
		logs.Err(err)
		return
	}
	loadRelay()
	go common.NewMonitor(false, "udp", fmt.Sprintf(":%d", common.Conf.ServerPort), handlerMsg).Start()
	go common.NewMonitor(false, "tcp", fmt.Sprintf(":%d", common.Conf.ServerPort-1), handlerMsg).Start()
	go common.NewMonitor(false, "udp", fmt.Sprintf(":%d", common.Conf.RegPort), handlerSyncMsg).Start()
	go common.NewMonitor(false, "tcp", fmt.Sprintf(":%d", common.Conf.ServerPort), handlerMsg).Start()
	api_server.Start()
}

// 黑名单检测
func blacklistDetection(id string, addr *common.Addr) bool {
	in := common.InList(addr.GetIP())
	if common.Conf.WhiteList && in {
		return false
	}
	if !common.Conf.WhiteList && in {
		return true
	}
	return false
}
func loadRelay() {
	online, err := dataSever.GetRelayAllOnline()
	if err == nil && len(online) > 0 {
		rendezvousServers = make([]string, len(online))
		for _, relay := range online {
			rendezvousServers = append(rendezvousServers, fmt.Sprintf("%s:%d", relay.IP, relay.Port))
		}
	}
}

// 获取一个中继服务器
func getRelay() string {
	online, _ := dataSever.GetRelayAllOnline()
	if len(online) != len(rendezvousServers) {
		loadRelay()
	}
	if len(online) == 0 {
		logs.Err("NoRegisteredRelayServer")
		return ""
	}
	var Relay *model.Relay
	for _, relay := range online {
		if relay.Cpu < 60 && (float64(relay.Upload)-relay.NetFlow) > (float64(relay.Upload)*0.1) {
			Relay = relay
			break
		}
	}
	if Relay == nil {
		Relay = online[rand.Intn(len(online))]
	}
	return fmt.Sprintf("%s:%d", Relay.IP, Relay.Port)
}
