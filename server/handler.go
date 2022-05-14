package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type ringMsg struct {
	ID      string //消息发送者ID
	Type    string
	TimeOut uint32
	InsTime time.Time
	Val     interface{}
}

func getMsgForm(id, Type string, timeOut uint) interface{} {
	if timeOut == 0 {
		timeOut = 3
	}
	after := time.After(time.Second * time.Duration(timeOut))
	for {
		select {
		case <-after:
			return nil
		default:
			next := r.Next()
			val := next.Val()
			now := time.Now()
			if val != nil {
				if v, ok := val.(*ringMsg); ok {
					if now.Add(time.Second * time.Duration(v.TimeOut)).Before(now) {
						next.Set(nil)
					} else if v.ID == id && v.Type == Type {
						next.Set(nil)
						return v.Val
					}
				}
			}
		}
	}
}
func handlerMsg(msg []byte, writer *common.Writer) {
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil {
		logs.Err(err)
	}
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessagePunchHoleRequest:
		//打洞
		HoleRequest := message.GetPunchHoleRequest()
		if HoleRequest == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessagePunchHoleRequest(HoleRequest))
	case model_proto.TypeRendezvousMessageRegisterPk:
		//注册公钥
		RegisterPk := message.GetRegisterPk()
		if RegisterPk == nil {
			return
		}
		pk := RendezvousMessageRegisterPk(RegisterPk)
		response = model_proto.NewRendezvousMessage(pk)
		if pk.Result == model_proto.RegisterPkResponse_OK {
			if _, ok := connPeerMap[RegisterPk.GetId()]; !ok {
				connPeerMap[RegisterPk.GetId()] = writer
			}
		}
	case model_proto.TypeRendezvousMessageRegisterPeer:
		//注册id
		RegisterPeer := message.GetRegisterPeer()
		if RegisterPeer == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessageRegisterPeer(RegisterPeer))
	case model_proto.TypeRendezvousMessageSoftwareUpdate:
		//软件更新
		SoftwareUpdate := message.GetSoftwareUpdate()
		if SoftwareUpdate == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessageSoftwareUpdate(SoftwareUpdate))
	case model_proto.TypeRendezvousMessageTestNatRequest:
		//网络类型测试
		TestNatRequest := message.GetTestNatRequest()
		if TestNatRequest == nil {
			return
		}
		request := RendezvousMessageTestNatRequest(TestNatRequest)
		str := writer.GetAddrStr()
		split := strings.Split(str, ":")
		parseUint, _ := strconv.ParseUint(split[1], 10, 32)
		request.Port = int32(parseUint)
		response = model_proto.NewRendezvousMessage(request)
	case model_proto.TypeRendezvousMessageLocalAddr:
		//本地地址返回
		LocalAddr := message.GetLocalAddr()
		if LocalAddr == nil {
			return
		}
		RendezvousMessageLocalAddr(LocalAddr)
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
