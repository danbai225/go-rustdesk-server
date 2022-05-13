package my_bytes

import (
	"bytes"
	"encoding/binary"
	logs "github.com/danbai225/go-logs"
	"github.com/shabbyrobe/go-num"
	"strconv"
	"strings"
	"time"
)

func encodeAddr(addr string) (b []byte) {
	b = make([]byte, 0)
	tm := num.U128From32(uint32(time.Now().UnixMicro()))
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return
	}
	pInt, _ := strconv.ParseUint(split[1], 10, 32)
	for _, s := range strings.Split(split[0], ".") {
		parseInt, _ := strconv.ParseUint(s, 10, 8)
		b = append(b, byte(parseInt))
	}
	var y uint32
	_ = binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &y)
	ip := num.U128From32(y)
	port := num.U128From32(uint32(pInt))
	logs.Info(port)
	ip.Add(tm)
	//v := ((ip + tm) << 49) | (tm << 17) | (port + (tm & 0xFFFF))
	return nil
}
func decodeAddr([]byte) string {
	return ""
}
