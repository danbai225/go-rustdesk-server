package server

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/container/gring"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
)

var dataSever data_server.DataSever
var queue = gqueue.New()
var r = gring.New(32, true)
var rendezvousServers = []string{"1.14.47.89"}
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
	go common.NewMonitor("udp", fmt.Sprintf(":%d", common.Conf.ServerPort), handlerMsg).Start()
	go common.NewMonitor("tcp", fmt.Sprintf(":%d", common.Conf.ServerPort-1), handlerMsg).Start()
	go common.NewMonitor("udp", fmt.Sprintf(":%d", common.Conf.RegPort), handlerSyncMsg).Start()
	common.NewMonitor("tcp", fmt.Sprintf(":%d", common.Conf.ServerPort), handlerMsg).Start()
}

//黑名单检测
func blacklistDetection(id string, addr *common.Addr) bool {
	if common.InList(addr.GetIP()) && !common.Conf.WhiteList {
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
