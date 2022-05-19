package model_proto

import (
	"reflect"
	"strings"
)

var (
	TypeMessageMessageCliprdr       = reflect.TypeOf(&Message_Cliprdr{}).String()
	TypeMessageMessageLoginResponse = reflect.TypeOf(&Message_LoginResponse{}).String()

	TypeRendezvousMessagePunchHoleRequest = reflect.TypeOf(&RendezvousMessage_PunchHoleRequest{}).String()
	TypeRendezvousMessageRegisterPk       = reflect.TypeOf(&RendezvousMessage_RegisterPk{}).String()
	TypeRendezvousMessageRegisterPeer     = reflect.TypeOf(&RendezvousMessage_RegisterPeer{}).String()
	TypeRendezvousMessageSoftwareUpdate   = reflect.TypeOf(&RendezvousMessage_SoftwareUpdate{}).String()
	TypeRendezvousMessageTestNatRequest   = reflect.TypeOf(&RendezvousMessage_TestNatRequest{}).String()
	TypeRendezvousMessageLocalAddr        = reflect.TypeOf(&RendezvousMessage_LocalAddr{}).String()
	TypeRendezvousMessageRequestRelay     = reflect.TypeOf(&RendezvousMessage_RequestRelay{}).String()
	TypeRendezvousMessageRelayResponse    = reflect.TypeOf(&RendezvousMessage_RelayResponse{}).String()
	TypeRendezvousMessagePunchHoleSent    = reflect.TypeOf(&RendezvousMessage_PunchHoleSent{}).String()
	TypeRendezvousMessageConfigureUpdate  = reflect.TypeOf(&RendezvousMessage_ConfigureUpdate{}).String()
)

var typeMap map[string]reflect.Type

func init() {
	typeMap = map[string]reflect.Type{
		"RendezvousMessage_RegisterPkResponse":   reflect.TypeOf(&RendezvousMessage_RegisterPkResponse{}).Elem(),
		"RendezvousMessage_RegisterPeerResponse": reflect.TypeOf(&RendezvousMessage_RegisterPeerResponse{}).Elem(),
		"RendezvousMessage_SoftwareUpdate":       reflect.TypeOf(&RendezvousMessage_SoftwareUpdate{}).Elem(),
		"RendezvousMessage_FetchLocalAddr":       reflect.TypeOf(&RendezvousMessage_FetchLocalAddr{}).Elem(),
		"RendezvousMessage_PunchHoleResponse":    reflect.TypeOf(&RendezvousMessage_PunchHoleResponse{}).Elem(),
		"RendezvousMessage_TestNatResponse":      reflect.TypeOf(&RendezvousMessage_TestNatResponse{}).Elem(),
		"RendezvousMessage_RequestRelay":         reflect.TypeOf(&RendezvousMessage_RequestRelay{}).Elem(),
		"RendezvousMessage_RelayResponse":        reflect.TypeOf(&RendezvousMessage_RelayResponse{}).Elem(),
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
