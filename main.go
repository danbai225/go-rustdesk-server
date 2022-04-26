package main

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
	"google.golang.org/protobuf/proto"
	"net"
)

func main() {
	go func() {
		listen, err := net.Listen("tcp", ":21116")
		if err != nil {
			logs.Err(err)
			return
		}
		accept, err := listen.Accept()
		if err != nil {
			logs.Err(err)
			return
		}
		bytes := make([]byte, 1024)
		read, err := accept.Read(bytes)
		if err != nil {
			logs.Err(err)
			return
		}
		logs.Info(read)
		request := &model_proto.Message{}
		logs.Info(bytes[:read])

		err = proto.Unmarshal(bytes[:read], request)
		if err != nil {
			logs.Err(err)
			return
		}

	}()
	listen, err := net.Listen("tcp", ":21117")
	if err != nil {
		logs.Err(err)
		return
	}
	accept, err := listen.Accept()
	if err != nil {
		logs.Err(err)
		return
	}
	bytes := make([]byte, 1024)
	read, err := accept.Read(bytes)
	if err != nil {
		logs.Err(err)
		return
	}
	request := model_proto.PunchHoleRequest{}
	err = proto.Unmarshal(bytes[:read], &request)
	if err != nil {
		logs.Err(err)
		return
	}
	logs.Info(request.Id)
}
