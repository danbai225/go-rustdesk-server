package relay

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func Start() {
	go registered()
	go common.NewMonitor("udp", ":21117", handlerMsg).Start()
	common.NewMonitor("tcp", ":21117", handlerMsg).Start()
}
func registered() {
	//dial, err := net.Dial("tcp", ":21116")
	//if err != nil {
	//	logs.Err(err)
	//	return
	//}
	//dial.Write(model_proto.RegisterPk{})
}
func handlerMsg(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil {
		logs.Err(err)
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	default:
		logs.Info(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		marshal, err2 := proto.Marshal(response)
		if err2 != nil {
			logs.Err(err2)
			return
		}
		_, err2 = writer.Write(marshal)
		if err2 != nil {
			logs.Err(err2)
		}
	}
}
