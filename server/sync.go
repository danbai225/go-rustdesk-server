package server

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/model"
	"go-rustdesk-server/model/model_msg"
)

//同步服务
func handlerSyncMsg(msg []byte, writer *common.Writer) {
	if len(msg) == 0 {
		return
	}
	m := model_msg.Msg{}
	if m.MsgType == 0 {
		return
	}
	err := json.Unmarshal(msg, &m)
	if err != nil {
		logs.Err(err)
		return
	}
	switch m.MsgType {
	case model_msg.RegType:
		if m.RegMsg == nil {
			return
		}
		regRelay(m.RegMsg, writer)
	default:
		logs.Debug(m.MsgType)
	}
}

//regRelay
func regRelay(msg *model_msg.RegMsg, writer *common.Writer) {
	m := model_msg.Msg{
		Base:    model_msg.Base{MsgType: model_msg.RegRType},
		RegMsgR: &model_msg.RegMsgR{Err: ""},
	}
	defer func() {
		writer.SendJsonMsg(&m)
	}()
	if msg.Name == "" {
		m.RegMsgR.Err = model_msg.Err
		return
	}
	relay, err := dataSever.GetRelayByName(msg.Name)
	if err == nil {
		logs.Err(err)
		m.RegMsgR.Err = err.Error()
		return
	}
	newRelay := &model.Relay{
		Name:        msg.Name,
		Port:        msg.RelayPort,
		IP:          writer.GetAddr().GetIP(),
		Online:      true,
		LastRegTime: &m.Time,
	}
	if relay == nil {
		err = dataSever.AddRelay(newRelay)
		if err != nil {
			m.RegMsgR.Err = err.Error()
		}
		common.GetList()
		writer.SendJsonMsg(&model_msg.Msg{
			Base: model_msg.Base{MsgType: model_msg.RegRType},
			SyncList: &model_msg.SyncList{
				WhiteList: common.Conf.WhiteList,
				List:      common.GetList(),
			}})
		loadRelay()
	} else if relay.IP == writer.GetAddr().GetIP() {
		err = dataSever.UpdateRelay(newRelay)
		if err != nil {
			m.RegMsgR.Err = err.Error()
		}
	} else {
		m.RegMsgR.Err = model_msg.ExistName
	}
}
