package common

import (
	"crypto/ed25519"
	"encoding/base64"
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"os"
)

var pk []byte
var pkStr string
var sk []byte

func genKey() error {
	key, privateKey, _ := ed25519.GenerateKey(nil)
	err := ioutil.WriteFile(keyPath, privateKey, os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(keyPath+".pub", key, os.ModePerm)
	return err
}
func LoadKey() {
	if !Exists(keyPath) {
		os.MkdirAll("key", os.ModePerm)
		err := genKey()
		if err != nil {
			logs.Err("gen key err :", err)
			return
		}
	}
	var err error
	sk, err = ioutil.ReadFile(keyPath)
	if err != nil {
		logs.Err("open key err:", err)
		return
	}
	pk, err = ioutil.ReadFile(keyPath + ".pub")
	if err != nil {
		logs.Err("open key err:", err)
		return
	}
	pkStr = base64.StdEncoding.EncodeToString(pk)
	logs.Info("key=", pkStr)
}
func Sign(data []byte) []byte {
	return append(ed25519.Sign(sk, data), data...)
}
func GetSignPK(version, id string, peerPK []byte) []byte {
	bytes := make([]byte, 0)
	if version == "" || id == "" {
		return bytes
	}
	logs.Debug(id, peerPK)
	marshal, _ := proto.Marshal(&model_proto.IdPk{
		Id: id,
		Pk: peerPK,
	})
	return Sign(marshal)
}
