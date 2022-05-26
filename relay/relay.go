package relay

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/goccy/go-json"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model/model_msg"
	"go-rustdesk-server/model/model_proto"
	"go-rustdesk-server/my_bytes"
	"google.golang.org/protobuf/proto"
	"net"
	"os"
	"reflect"
	"time"
)

func Start() {
	go regRelay()
	common.NewMonitor(true, "tcp", fmt.Sprintf(":%d", common.Conf.RelayPort), handlerMsg).Start()
}

func handlerMsg(msg []byte, writer *common.Writer) {
	if blacklistDetection(writer.GetAddr()) {
		writer.Close()
		return
	}
	message := model_proto.RendezvousMessage{}
	err := proto.Unmarshal(msg, &message)
	if err != nil || message.Union == nil {
		if err != nil {
			logs.Err(err)
		}
		return
	}
	logs.Debug(writer.Type(), writer.GetAddrStr(), reflect.TypeOf(message.Union).String())
	switch reflect.TypeOf(message.Union).String() {
	case model_proto.TypeRendezvousMessageRequestRelay:
		RequestRelay := message.GetRequestRelay()
		if RequestRelay == nil {
			return
		}
		uuid := RequestRelay.GetUuid()
		logs.Debug(RequestRelay.Id, RequestRelay.RelayServer, RequestRelay.Token, RequestRelay.Secure, my_bytes.DecodeAddr(RequestRelay.SocketAddr))
		if uuid != "" {
			w, err1 := common.GetWriter(uuid, common.TCP)
			if err1 != nil {
				writer.SetKey(uuid)
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

//黑名单检测
func blacklistDetection(addr *common.Addr) bool {
	in := common.InList(addr.GetIP())
	if common.Conf.WhiteList && in {
		return false
	}
	if !common.Conf.WhiteList && in {
		return true
	}
	return false
}
func regRelay() {
	var dial net.Conn
	var err error
	var read int
	go func() {
		for {
			time.Sleep(time.Second * 1)
			if dial != nil {
				marshal, _ := json.Marshal(model_msg.Msg{
					Base: model_msg.Base{
						MsgType: model_msg.RegType,
					},
					RegMsg: &model_msg.RegMsg{Name: common.Conf.RelayName, Time: time.Now(), RelayPort: common.Conf.RelayPort},
				})
				_, err1 := dial.Write(marshal)
				if err1 != nil {
					logs.Err(err1)
				}
			}
		}
	}()
	for {
		dial, err = net.Dial("udp", common.Conf.RegServer)
		if err != nil {
			logs.Err(err)
		}
		bytes := make([]byte, 1024)
		for err == nil {
			read, err = dial.Read(bytes)
			if err == nil {
				bs := make([]byte, read)
				copy(bs, bytes[:read])
				go handlerSyncMsg(bs, dial)
			}
		}
		time.Sleep(time.Second * 15)
	}
}
func handlerSyncMsg(msg []byte, writer net.Conn) {
	if len(msg) == 0 {
		return
	}
	m := model_msg.Msg{}
	err := json.Unmarshal(msg, &m)
	if err != nil {
		logs.Err(err)
		return
	}
	if m.MsgType == 0 {
		return
	}
	switch m.MsgType {
	case model_msg.RegRType:
		if m.RegMsgR == nil {
			return
		}
		switch m.RegMsgR.Err {
		case model_msg.ExistName:
			logs.Err(m.RegMsgR.Err)
			os.Exit(1)
		default:
			//logs.Debug("注册成功")
		}
	case model_msg.SyncListType:
		if m.SyncList == nil {
			return
		}
		common.UpDataList(m.SyncList.WhiteList, m.SyncList.List)
	default:
		logs.Debug(m.MsgType)
	}
}
