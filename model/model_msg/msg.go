package model_msg

import "time"

const (
	RegType      = 10000
	RegRType     = 20000
	SyncListType = 10001
)
const (
	Err       = "Err"
	ExistName = "ExistName"
)

type Base struct {
	MsgType uint32 `json:"msg_type"`
}
type RegMsg struct {
	Name      string    `json:"name"`
	Time      time.Time `json:"time"`
	RelayPort uint16    `json:"relay_port"`
}
type RegMsgR struct {
	Err string `json:"err"`
}
type SyncList struct {
	WhiteList bool     `json:"white_list"`
	List      []string `json:"list"`
}
type Msg struct {
	Base
	*RegMsg
	*RegMsgR
	*SyncList
}
