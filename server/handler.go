package server

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"reflect"
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
	if blacklistDetection("", writer.GetAddr()) {
		return
	}
	logs.Debug(writer.Type(), writer.GetAddrStr(), reflect.TypeOf(message.Union).String())
	var response proto.Message
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessagePunchHoleRequest:
		//打洞
		HoleRequest := message.GetPunchHoleRequest()
		if HoleRequest == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessagePunchHoleRequest(HoleRequest, writer))
	case model_proto.TypeRendezvousMessageRegisterPk:
		//注册公钥
		RegisterPk := message.GetRegisterPk()
		if RegisterPk == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessageRegisterPk(RegisterPk, writer))
	case model_proto.TypeRendezvousMessageRegisterPeer:
		//注册id
		RegisterPeer := message.GetRegisterPeer()
		if RegisterPeer == nil {
			return
		}
		peer := RendezvousMessageRegisterPeer(RegisterPeer, writer)
		response = model_proto.NewRendezvousMessage(peer)
		ConfigureUpdate(writer)
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
		response = model_proto.NewRendezvousMessage(RendezvousMessageTestNatRequest(TestNatRequest, writer))
	case model_proto.TypeRendezvousMessageLocalAddr:
		//本地地址返回
		LocalAddr := message.GetLocalAddr()
		if LocalAddr == nil {
			return
		}
		RendezvousMessageLocalAddr(LocalAddr, writer)
	case model_proto.TypeRendezvousMessageRequestRelay:
		//请求继中
		RequestRelay := message.GetRequestRelay()
		if RequestRelay == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessageRequestRelay(RequestRelay))
	case model_proto.TypeRendezvousMessageRelayResponse:
		//请求继中
		RelayResponse := message.GetRelayResponse()
		if RelayResponse == nil {
			return
		}
		RendezvousMessageRelayResponse(RelayResponse)
	case model_proto.TypeRendezvousMessagePunchHoleSent:
		//请求打洞
		PunchHoleSent := message.GetPunchHoleSent()
		if PunchHoleSent == nil {
			return
		}
		response = model_proto.NewRendezvousMessage(RendezvousMessagePunchHoleSent(PunchHoleSent, writer))
	case model_proto.TypeRendezvousMessageConfigureUpdate:
		//配置更新
		ConfigUpdate := message.GetConfigureUpdate()
		if ConfigUpdate == nil {
			return
		}
		RendezvousMessageConfigureUpdate(ConfigUpdate)
	default:
		logs.Debug(reflect.TypeOf(message.Union).String())
	}
	if response != nil {
		writer.SendMsg(response)
	}
}
