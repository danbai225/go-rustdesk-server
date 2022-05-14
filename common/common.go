package common

import (
	"errors"
	"fmt"
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
	call    func(msg []byte, writer *Writer)
}
type Writer struct {
	Type  string
	tConn net.Conn
	uConn *net.UDPConn
	addr  *net.UDPAddr
}

func (w *Writer) Write(p []byte) (n int, err error) {
	switch w.Type {
	case "udp", "UDP":
		return w.uConn.WriteToUDP(p, w.addr)
	case "TCP", "tcp":
		encoder, err := my_bytes.Encoder(p)
		if err != nil {
			return 0, err
		}
		return w.tConn.Write(encoder)
	}
	return 0, errors.New("type Err")
}
func (w *Writer) GetAddrStr() string {
	switch w.Type {
	case "udp", "UDP":
		return w.addr.String()
	case "TCP", "tcp":
		return w.tConn.RemoteAddr().String()
	}
	return ""
}

var myWriterMap = map[string]*Writer{}

func NewMonitor(network, addr string, call func(msg []byte, writer *Writer)) *monitor {
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
		var writer *Writer
		var ok bool
		k := fmt.Sprint("udp", addr.String())
		if writer, ok = myWriterMap[k]; !ok {
			writer = &Writer{
				Type:  "udp",
				uConn: m.conn,
				addr:  addr,
			}
			myWriterMap[k] = writer
		}
		m.call(temp, writer)
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
	data, err = my_bytes.Decode(data)
	if err != nil {
		logs.Err(err)
		return
	}
	var writer *Writer
	var ok bool
	k := fmt.Sprint("tcp", conn.RemoteAddr().String())
	if writer, ok = myWriterMap[k]; !ok {
		writer = &Writer{
			Type:  "tcp",
			tConn: conn,
		}
		myWriterMap[k] = writer
	}
	m.call(data, writer)
}
