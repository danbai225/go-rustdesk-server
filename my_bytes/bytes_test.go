package my_bytes

import (
	"bytes"
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"testing"
)

var testData = []byte{100, 66, 23, 10, 14, 49, 50, 51, 49, 50, 51, 113, 113, 113, 113, 113, 113, 113, 113, 16, 2, 26, 3, 49, 50, 51}

func TestDecode(t *testing.T) {
	data, err2 := Decode(testData)
	if err2 != nil {
		logs.Err(err2)
		return
	}
	msg := &model_proto.Message{}
	err := proto.Unmarshal(data, msg)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(msg.Union.(*model_proto.Message_LoginResponse).LoginResponse.String())
	encoder, err2 := Encoder(data)
	logs.Info(encoder, err2)
}
func tt(data []byte) {
	data[0] = byte(99)
}
func TestEncoder(t *testing.T) {
	b := bytes.NewBufferString("")
	for b.Len() <= 0x3a {
		b.WriteString("test")
	}
	cp := make([]byte, b.Len())
	copy(cp, b.Bytes())
	encode, err := Encoder(cp)
	if err != nil {
		t.Error(err)
	}
	data, err := Decode(encode)
	if err != nil {
		t.Error(err)
	}
	if b.String() != string(data) {
		t.Error("!=")
	}
	//0x3FFF
	b = bytes.NewBufferString("")
	for b.Len() <= 0x3FFF {
		b.WriteString("test")
	}
	cp = make([]byte, b.Len())
	copy(cp, b.Bytes())
	encode, err = Encoder(cp)
	if err != nil {
		t.Error(err)
	}
	data, err = Decode(encode)
	if err != nil {
		t.Error(err)
	}
	if b.String() != string(data) {
		t.Error("!=")
	}
	logs.Info(b.String())
	logs.Info(string(data))
}
