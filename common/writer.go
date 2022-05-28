package common

import (
	"encoding/json"
	"errors"
	"fmt"
	logs "github.com/danbai225/go-logs"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"go-rustdesk-server/model/model_msg"
	"go-rustdesk-server/model/model_proto"
	"go-rustdesk-server/my_bytes"
	"google.golang.org/protobuf/proto"
	"net"
	"strconv"
	"strings"
	"time"
)

var ctx = gctx.New()
var cache = gcache.New()
var cacheTimeOut = time.Second * 60

//func init() {
//	go func() {
//		for true {
//			c, _ := cache.Data(ctx)
//			strings := make([]string, 0)
//			for k, _ := range c {
//				strings = append(strings, k.(string))
//			}
//			sort.Strings(strings)
//			time.Sleep(time.Second)
//			logs.Info(strings)
//		}
//	}()
//}

type Writer struct {
	key   string
	_type string
	tConn net.Conn
	uConn *net.UDPConn
	addr  *net.UDPAddr
	loop  bool
}
type Addr struct {
	ip   string
	port uint32
}

func (a *Addr) GetIP() string {
	return a.ip
}
func (a *Addr) GetPort() uint32 {
	return a.port
}
func (a *Addr) Parsing(addr string) {
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return
	}
	a.ip = split[0]
	p, _ := strconv.ParseUint(split[1], 10, 32)
	a.port = uint32(p)
}
func (w *Writer) Type() string {
	return w._type
}
func (w *Writer) Write(p []byte) (n int, err error) {
	switch w._type {
	case udp:
		if w.uConn == nil {
			return 0, errors.New("uConn==nil")
		}
		return w.uConn.WriteToUDP(p, w.addr)
	case tcp:
		if w.tConn == nil {
			return 0, errors.New("tConn==nil")
		}
		encoder, err := my_bytes.Encoder(p)
		if err != nil {
			return 0, err
		}
		return w.tConn.Write(encoder)
	}
	return 0, errors.New("type Err")
}
func (w *Writer) WriteToAddr(p []byte, addr string) (n int, err error) {
	switch w._type {
	case udp:
		if w.uConn == nil {
			return 0, errors.New("uConn==nil")
		}
		udpAddr, err1 := net.ResolveUDPAddr(udp, addr)
		if err1 != nil {
			return 0, err1
		}
		return w.uConn.WriteToUDP(p, udpAddr)
	case tcp:
		//if w.tConn == nil {
		//	return 0, errors.New("tConn==nil")
		//}
		//encoder, err := my_bytes.Encoder(p)
		//if err != nil {
		//	return 0, err
		//}
		//return w.tConn.Write(encoder)
		return 0, errors.New("unrealized")
	}
	return 0, errors.New("type Err")
}
func (w *Writer) GetAddrStr() string {
	switch w._type {
	case udp:
		return w.addr.String()
	case tcp:
		return w.tConn.RemoteAddr().String()
	}
	return ""
}
func (w *Writer) GetAddr() *Addr {
	addr := ""
	switch w._type {
	case udp:
		addr = w.addr.String()
	case tcp:
		addr = w.tConn.RemoteAddr().String()
	}
	a := &Addr{}
	a.Parsing(addr)
	return a
}
func (w *Writer) SetKey(key string) {
	mk := ""
	switch w._type {
	case udp:
		mk = udp + key
	case tcp:
		mk = tcp + key
	}
	_ = cache.Set(ctx, mk, w, cacheTimeOut)
	w.key = key
}
func (w *Writer) setAddr(addr string) {
	mk := ""
	switch w._type {
	case udp:
		mk = udp + addr
	case tcp:
		mk = tcp + addr
	}
	_ = cache.Set(ctx, mk, w, time.Second*60)
}
func (w *Writer) remove() {
	mk := ""
	switch w._type {
	case udp:
		mk = udp + w.addr.String()
	case tcp:
		mk = tcp + w.tConn.RemoteAddr().String()
	}
	_, _ = cache.Remove(ctx, mk)
	if w.key != "" {
		mk = udp + w.key
		_, _ = cache.Remove(ctx, mk)
		mk = tcp + w.key
		_, _ = cache.Remove(ctx, mk)
	}
}
func (w *Writer) Copy(dst *Writer) {
	if w._type != tcp || dst == nil || dst.tConn == nil {
		return
	}
	//go io.Copy(dst.tConn, w.tConn)
	//io.Copy(w.tConn, dst.tConn)
	a := func(d, s net.Conn) {
		bytes := make([]byte, 10240)
		for {
			read, err := s.Read(bytes)
			if err != nil {
				logs.Err(err)
				return
			} else {
				message := model_proto.Message{}
				err = proto.Unmarshal(bytes, &message)
				if err == nil {
					logs.Info(message.Union)
				}
			}
			write, err := d.Write(bytes[:read])
			if err != nil || write != read {
				logs.Err(err)
				return
			}
			logs.Info(write)
		}
	}
	go a(dst.tConn, w.tConn)
	a(w.tConn, dst.tConn)
}
func (w *Writer) SendMsg(message proto.Message) {
	if message == nil {
		return
	}
	marshal, err2 := proto.Marshal(message)
	if err2 != nil {
		logs.Err(err2)
		return
	}
	_, err2 = w.Write(marshal)
	if err2 != nil {
		logs.Err(err2)
	}
}
func (w *Writer) SendJsonMsg(message *model_msg.Msg) {
	if message == nil {
		return
	}
	marshal, err2 := json.Marshal(message)
	if err2 != nil {
		logs.Err(err2)
		return
	}
	_, err2 = w.Write(marshal)
	if err2 != nil {
		logs.Err(err2)
	}
}
func (w *Writer) Close() {
	if w._type == tcp {
		_ = w.tConn.Close()
	}
	w.remove()
}
func GetWriter(key, _type string) (*Writer, error) {
	mk := ""
	switch _type {
	case udp:
		mk = fmt.Sprint(udp, key)
	case tcp:
		mk = fmt.Sprint(tcp, key)
	}
	get, _ := cache.Get(ctx, mk)
	if get != nil {
		if v, ok := get.Val().(*Writer); ok {
			return v, nil
		}
	}
	return nil, errors.New("OFFLINE")
}
func addWriter(key, _type string, w *Writer) {
	mk := ""
	t := cacheTimeOut
	switch _type {
	case udp:
		mk = fmt.Sprint(udp, key)
	case tcp:
		mk = fmt.Sprint(tcp, key)
		t = 0
	}
	err := cache.Set(ctx, mk, w, t)
	if err != nil {
		logs.Err(err)
	}
}
func RemoveWriter(w *Writer) {
	w.remove()
}
