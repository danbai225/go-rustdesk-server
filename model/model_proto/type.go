package model_proto

import (
	"reflect"
	"strings"
)

var (
	TypeMessageMessageCliprdr       = reflect.TypeOf(&Message_Cliprdr{}).String()
	TypeMessageMessageLoginResponse = reflect.TypeOf(&Message_LoginResponse{}).String()

	TypeRendezvousMessageRegisterPk     = reflect.TypeOf(&RendezvousMessage_RegisterPk{}).String()
	TypeRendezvousMessageRegisterPeer   = reflect.TypeOf(&RendezvousMessage_RegisterPeer{}).String()
	TypeRendezvousMessageSoftwareUpdate = reflect.TypeOf(&RendezvousMessage_SoftwareUpdate{}).String()
)

var typeMap map[string]reflect.Type

func init() {
	typeMap = map[string]reflect.Type{
		"RendezvousMessage_RegisterPkResponse":   reflect.TypeOf(&RendezvousMessage_RegisterPkResponse{}).Elem(),
		"RendezvousMessage_RegisterPeerResponse": reflect.TypeOf(&RendezvousMessage_RegisterPeerResponse{}).Elem(),
		"RendezvousMessage_SoftwareUpdate":       reflect.TypeOf(&RendezvousMessage_SoftwareUpdate{}).Elem(),
	}
}

func getTypeByName(name string) reflect.Type {
	if v, has := typeMap[name]; has {
		return v
	} else if strings.Contains(name, "*") {
		return getTypeByName(strings.ReplaceAll(name, "*", ""))
	}
	return nil
}
