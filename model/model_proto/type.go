package model_proto

import (
	"reflect"
	"strings"
)

var (
	TypeRendezvousMessageRegisterPk   = reflect.TypeOf(&RendezvousMessage_RegisterPk{}).String()
	TypeRendezvousMessageRegisterPeer = reflect.TypeOf(&RendezvousMessage_RegisterPeer{}).String()
)

var typeMap map[string]reflect.Type

func init() {
	typeMap = map[string]reflect.Type{
		"RendezvousMessage_RegisterPkResponse":   reflect.TypeOf(&RendezvousMessage_RegisterPkResponse{}).Elem(),
		"RendezvousMessage_RegisterPeerResponse": reflect.TypeOf(&RendezvousMessage_RegisterPeerResponse{}).Elem(),
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
