package server

import (
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/v2/container/gqueue"
	"github.com/gogf/gf/v2/container/gring"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
)

var dataSever data_server.DataSever
var ctx = gctx.New()
var cache = gcache.New()
var queue = gqueue.New()
var r = gring.New(10, true)
var connPeerMap = map[string]*common.Writer{}

func Start() {
	var err error
	dataSever, err = data_server.GetDataSever()
	if err != nil {
		logs.Err(err)
		return
	}
	go common.NewMonitor("udp", ":21116", handlerMsg).Start()
	common.NewMonitor("tcp", ":21116", handlerMsg).Start()
}
