package common

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/my_bytes"
	"go.uber.org/zap/buffer"
	"io"
	"net"
)

type monitor struct {
	network string
	addr    string
	listen  net.Listener
	conn    *net.UDPConn
	call    func(msg []byte, write func(data []byte) error, conn net.Conn)
}

func NewMonitor(network, addr string, call func(msg []byte, write func(data []byte) error, conn net.Conn)) *monitor {
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
		addr, err1 := net.ResolveUDPAddr(m.network, m.addr)
		if err1 != nil {
			logs.Err(err1)
		}
		m.conn, err = net.ListenUDP(m.network, addr)
		if err != nil {
			logs.Err(err)
		}
		m.readUdp()
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
	realLength := uint(0)
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
		if m.network != "udp" {
			_, _ = bytes.Write(temp)
			if realLength == 0 && bytes.Len() > 0 {
				_, realLength, err = my_bytes.DecodeHead(bytes.Bytes())
				if err != nil {
					logs.Err(err)
					return
				}
			}
			//长度不匹配继续读
			if int(realLength) < bytes.Len() {
				logs.Err("int(length) <bytes.Len()")
				return
			} else if int(realLength) == bytes.Len() {
				//解析协议
				cp := make([]byte, bytes.Len())
				copy(cp, bytes.Bytes())
				go m.processMessageData(cp[:bytes.Len()], conn)
				bytes.Reset()
			}
		}
	}
}
func (m *monitor) readUdp() {
	for {
		temp := make([]byte, 1024)
		readLen, addr, err := m.conn.ReadFromUDP(temp)
		if err != nil && err != io.EOF {
			logs.Err(err)
			return
		}
		if readLen == 0 {
			continue
		}
		temp = temp[:readLen]
		m.call(temp, func(data []byte) error {
			_, err2 := m.conn.WriteToUDP(data, addr)
			return err2
		}, nil)
	}
}

func (m *monitor) processMessageData(data []byte, conn net.Conn) {
	defer func() {
		err := recover()
		if err != nil {
			logs.Err(err)
		}
	}()
	var err error
	if m.network != "udp" {
		data, err = my_bytes.Decode(data)
		if err != nil {
			logs.Err(err)
			return
		}
	}
	m.call(data, func(data []byte) error {
		_, err2 := conn.Write(data)
		return err2
	}, conn)
}
