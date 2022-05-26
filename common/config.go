package common

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"os"
)

type Config struct {
	WhiteList  bool     `json:"whiteList"`
	IpList     []string `json:"ipList"`
	Debug      bool     `json:"debug"`
	RegServer  string   `json:"reg_server"`
	RelayName  string   `json:"relay_name"`
	RegPort    uint16   `json:"reg_port"`
	ServerPort uint16   `json:"server_port"`
	RelayPort  uint16   `json:"relay_port"`
}

var Conf *Config
var path string

func init() {
	path = "config.json"
	env, _ := os.LookupEnv("RUST_DESK_CONF_PATH")
	if env != "" {
		path = env
	}
	file, err := ioutil.ReadFile(path)
	if err != nil {
		logs.Err(err)
		return
	}
	Conf = &Config{}
	err = json.Unmarshal(file, Conf)
	if err != nil {
		logs.Err(err)
		return
	}
	loadList()
}
func Save() {
	marshal, _ := json.Marshal(Conf)
	err := ioutil.WriteFile(path, marshal, os.ModePerm)
	if err != nil {
		logs.Err(err)
	}
}
