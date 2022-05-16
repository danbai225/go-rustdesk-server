package common

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
	"go-rustdesk-server/my_bytes"
	"io"
	"net"
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
func (w *Writer) GetAddrStr() string {
	switch w._type {
	case udp:
		return w.addr.String()
	case tcp:
		return w.tConn.RemoteAddr().String()
	}
	return ""
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
	go io.Copy(dst.tConn, w.tConn)
	io.Copy(w.tConn, dst.tConn)
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
	cache.Set(ctx, mk, w, t)
}
