package common

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/model/model_proto"
	"go-rustdesk-server/my_bytes"
	"go.uber.org/zap/buffer"
	"google.golang.org/protobuf/proto"
	"io"
	"net"
)

type monitor struct {
	network string
	addr    string
	listen  net.Listener
	conn    net.Conn
	call    func(msg *model_proto.Message)
}

func NewMonitor(network, addr string, call func(msg *model_proto.Message)) *monitor {
	return &monitor{network: network, addr: addr, call: call}
}
func (m *monitor) Start() {
	defer func() {
		if m.listen != nil {
			_ = m.listen.Close()
		}
		if m.conn != nil {
			_ = m.conn.Close()
		}
	}()
	var err error
	if m.network == "udp" {
		m.conn, err = net.ListenUDP(m.network, &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 21116})
		m.accept(m.conn)
	} else {
		m.listen, err = net.Listen(m.network, m.addr)
		if err != nil {
			logs.Err(err)
			return
		}
		for {
			conn, err2 := m.listen.Accept()
			if err2 != nil {
				logs.Err(err2)
			} else {
				go m.accept(conn)
			}
		}
	}
}
func (m *monitor) accept(conn net.Conn) {
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()
	bytes := buffer.NewPool().Get()
	for {
		temp := make([]byte, 1024)
		readLen, err := conn.Read(temp)
		if err != nil && err != io.EOF {
			logs.Err(err)
			return
		}
		if readLen == 0 {
			continue
		}
		temp = temp[:readLen]
		_, _ = bytes.Write(temp)
		length, err := my_bytes.DecodeHead(bytes.Bytes())
		if err != nil {
			logs.Err(err)
			return
		}
		//长度不匹配继续读
		if int(length) < bytes.Len() {
			logs.Err("int(length) <bytes.Len()")
			return
		} else if int(length) == bytes.Len() {
			//解析协议
			cp := make([]byte, bytes.Len())
			copy(cp, bytes.Bytes())
			go m.processMessage(cp)
			bytes.Reset()
		}
	}
}
func (m *monitor) processMessage(data []byte) {
	defer func() {
		err := recover()
		if err != nil {
			logs.Err(err)
		}
	}()
	msg := &model_proto.Message{}
	decode, err := my_bytes.Decode(data)
	if err != nil {
		logs.Err(err)
		return
	}
	err = proto.Unmarshal(decode, msg)
	if err != nil {
		logs.Err(err)
		return
	}
	m.call(msg)
}
