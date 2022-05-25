package main

import (
	"flag"
	"go-rustdesk-server/data_server"
	"go-rustdesk-server/relay"
	"go-rustdesk-server/server"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	go http.ListenAndServe("0.0.0.0:21117", nil)
	_relay := flag.Bool("relay", false, "run relay")
	_server := flag.Bool("server", true, "run server")
	flag.Parse()
	if *_relay {
		go relay.Start()
	}
	if *_server {
		go server.Start()
	}
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		sever, err := data_server.GetDataSever()
		if err == nil {
			sever.Close()
		}
		done <- true
	}()
	<-done
}
