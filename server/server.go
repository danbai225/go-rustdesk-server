package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
)

func Start() {
	go common.NewMonitor("udp", ":21116", handlerMsg).Start()
	common.NewMonitor("tcp", ":21116", handlerMsg).Start()
}
func handlerMsg(message *model_proto.Message) {
	logs.Info(message.Union)
}
