package relay

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
)

func Start() {
	go common.NewMonitor("udp", ":21117", handlerMsg).Start()
	common.NewMonitor("tcp", ":21117", handlerMsg).Start()
}

var w *common.Writer

func handlerMsg(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil || message.Union == nil {
		if err != nil {
			logs.Err(err)
		}
		return
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRequestRelay:
		if w == nil {
			w = writer
		}
		if w != nil && w != writer {
			go w.Copy(writer)
		}
		//relay := message.GetRequestRelay()
		//if relay == nil {
		//	return
		//}
		//logs.Info(writer.GetAddrStr(), relay.GetId(), relay.GetConnType(), relay.GetSocketAddr(), relay.GetRelayServer(), relay.GetUuid())
	default:
		logs.Info(writer.GetAddrStr(), reflect.TypeOf(message.Union).String())
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
