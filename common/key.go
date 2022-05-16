package common

import (
	"crypto/ed25519"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"os"
)

var pk []byte
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
}
func Sign(data []byte) []byte {
	return ed25519.Sign(pk, data)
}
func Verify(data, sign []byte) bool {
	defer func() {
		err := recover()
		if err != nil {
			logs.Err(err)
		}
	}()
	return ed25519.Verify(pk, data, sign)
}
