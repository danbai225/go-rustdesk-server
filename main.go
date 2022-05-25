package main

import (
	"flag"
	"go-rustdesk-server/relay"
	"go-rustdesk-server/server"
)

func main() {
	_relay := flag.Bool("relay", false, "run relay")
	_server := flag.Bool("server", false, "run server")
	flag.Parse()
	if *_relay {
		relay.Start()
	} else if *_server {
		server.Start()
	}
}
