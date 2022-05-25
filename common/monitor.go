package common

import (
	logs "github.com/danbai225/go-logs"
	"go-rustdesk-server/my_bytes"
	"go.uber.org/zap/buffer"
	"net"
)

type monitor struct {
	network string
	addr    string
	listen  net.Listener
	conn    *net.UDPConn
	call    func(msg []byte, writer *Writer)
	relay   bool
}

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
	if m.network == udp {
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
				//logs.Info(conn.RemoteAddr().String())
				go m.accept(conn)
			}
		}
	}
}
func (m *monitor) accept(conn net.Conn) {
	writer := &Writer{
		_type: tcp,
		tConn: conn,
		loop:  true,
	}
	addWriter(conn.RemoteAddr().String(), tcp, writer)
	defer func() {
		if conn != nil {
			writer, _ := GetWriter(conn.RemoteAddr().String(), tcp)
			if writer != nil {
				writer.remove()
			}
			_ = conn.Close()
		}
	}()
	bytes := buffer.NewPool().Get()
	realLength := uint(0)
	temp := make([]byte, 1024)
	for writer.loop {
		readLen, err := conn.Read(temp)
		if err != nil {
			return
		}
		if readLen == 0 {
			continue
		}
		temp = temp[:readLen]
		_, _ = bytes.Write(temp)
		if realLength == 0 {
			_, realLength, err = my_bytes.DecodeHead(bytes.Bytes())
			if err != nil {
				//logs.Err(err)
				return
			}
		}
		if bytes.Len() >= int(realLength) {
			//解析协议
			cp := make([]byte, realLength)
			copy(cp, bytes.Bytes())
			if bytes.Len() != int(realLength) {
				bs := bytes.Bytes()[realLength:]
				bytes.Reset()
				bytes.Write(bs)
			}
			if m.relay && writer != nil {
				writer.loop = false
			}
			go m.processMessageData(cp, conn)
		}
	}
}
func (m *monitor) readUdp() {
	for {
		var writer *Writer
		temp := make([]byte, 1024)
		readLen, addr, err := m.conn.ReadFromUDP(temp)
		if err != nil {
			if writer != nil {
				writer.remove()
			}
			logs.Err(err)
			return
		}
		if readLen == 0 {
			continue
		}
		temp = temp[:readLen]
		writer, err = GetWriter(addr.String(), udp)
		if err != nil {
			writer = &Writer{
				_type: "udp",
				uConn: m.conn,
				addr:  addr,
			}
			addWriter(addr.String(), udp, writer)
		}
		m.call(temp, writer)
	}
}

func (m *monitor) processMessageData(data []byte, conn net.Conn) {
	var writer *Writer
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
	writer, err = GetWriter(conn.RemoteAddr().String(), tcp)
	if err != nil {
		writer = &Writer{
			_type: "tcp",
			tConn: conn,
		}
		addWriter(conn.RemoteAddr().String(), tcp, writer)
	}
	m.call(data, writer)
}
