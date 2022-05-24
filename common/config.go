package common

import (
	"encoding/json"
	logs "github.com/danbai225/go-logs"
	"io/ioutil"
	"os"
)

type Config struct {
	WhiteList bool     `json:"whiteList"`
	IpList    []string `json:"ipList"`
	Debug     bool     `json:"debug"`
}

var Conf *Config

func init() {
	path := "config.json"
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
