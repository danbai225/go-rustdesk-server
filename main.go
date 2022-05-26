package main

import (
	"flag"
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/common"
	"go-rustdesk-server/data_server"
	"go-rustdesk-server/relay"
	"go-rustdesk-server/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	_relay := flag.Bool("relay", common.Conf.RelayName != "", "run relay")
	_server := flag.Bool("server", true, "run server")
	flag.Parse()
	if *_relay {
		go relay.Start()
	}
	if *_server {
		go server.Start()
	}
	if common.Conf.Debug {
		logs.SetWriteLogs(logs.INFO | logs.ERR | logs.DEBUG)
	} else {
		logs.SetWriteLogs(logs.INFO | logs.ERR)
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		sever, err := data_server.GetDataSever()
		if err == nil {
			err = sever.Close()
			if err != nil {
				logs.Err(err)
			}
		}
		done <- true
	}()
	<-done
}
