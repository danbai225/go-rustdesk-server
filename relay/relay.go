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

func handlerMsg(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil || message.Union == nil {
		if err != nil {
			logs.Err(err)
		}
		return
	}
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRequestRelay:
		RequestRelay := message.GetRequestRelay()
		if RequestRelay == nil {
			return
		}
		uuid := RequestRelay.GetUuid()
		if uuid != "" {
			w, err1 := common.GetWriter(uuid, common.TCP)
			if err1 != nil {
				writer.SetKey(RequestRelay.GetUuid())
			} else if w != nil {
				common.RemoveWriter(writer)
				common.RemoveWriter(w)
				go writer.Copy(w)
			}
		}
	default:
		logs.Info(writer.GetAddrStr(), reflect.TypeOf(message.Union).String())
	}
}
