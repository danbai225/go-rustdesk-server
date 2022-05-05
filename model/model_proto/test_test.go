package model_proto

import (
	logs "github.com/danbai225/go-logs"
	"google.golang.org/protobuf/proto"
	"reflect"
	"strings"
	"testing"
)

func TestOK(t *testing.T) {
	message := RendezvousMessage{}
	message.Union = &RendezvousMessage_RegisterPkResponse{
		RegisterPkResponse: &RegisterPkResponse{
			Result: 0,
		},
	}
	marshal, err := proto.Marshal(&message)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(marshal)
	message2 := RendezvousMessage{}
	err = proto.Unmarshal(marshal, &message2)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(message2.Union)
}
func TestNewMsg(t *testing.T) {
	response := &RegisterPkResponse{
		Result: 0,
	}
	tt := reflect.TypeOf(response)
	ts := strings.ReplaceAll(strings.ReplaceAll(tt.String(), "*", ""), "model_proto.", "")
	ty := getTypeByName("RendezvousMessage_" + ts)
	if ty == nil {
		return
	}
	newStruc := reflect.New(ty)
	f := newStruc.Elem().FieldByName(ts)
	f.Set(reflect.ValueOf(response))
	logs.Info(newStruc.Elem())

	message := NewRendezvousMessage(response)
	logs.Info(message)
}
