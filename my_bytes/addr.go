package my_bytes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/shabbyrobe/go-num"
	"math"
	"strconv"
	"strings"
	"time"
)

func EncodeAddr(addr string) (bs []byte) {
	bs = make([]byte, 0)
	tm := num.U128From32(uint32(time.Now().UnixMicro()))
	split := strings.Split(addr, ":")
	if len(split) != 2 {
		return
	}
	pInt, _ := strconv.ParseUint(split[1], 10, 32)
	for _, s := range strings.Split(split[0], ".") {
		parseInt, _ := strconv.ParseUint(s, 10, 8)
		bs = append(bs, byte(parseInt))
	}
	var y uint32
	_ = binary.Read(bytes.NewBuffer(bs), binary.LittleEndian, &y)
	ip := num.U128From32(y)
	port := num.U128From32(uint32(pInt))
	//v := ((ip + tm) << 49) | (tm << 17) | (port + (tm & 0xFFFF))
	a := ip.Add(tm).Lsh(49)
	b := tm.Lsh(17)
	c := tm.And(num.U128From32(0xFFFF)).Add(port)
	d := a.Or(b).Or(c)
	bs = make([]byte, 16)
	d.PutLittleEndian(bs)
	nPadding := 0
	for i := 15; i >= 0; i-- {
		if bs[i] == 0 {
			nPadding += 1
		} else {
			break
		}
	}
	bs = bs[:(16 - nPadding)]
	return bs
}
func DecodeAddr(data []byte) string {
	bs := make([]byte, 16)
	copy(bs, data)
	number := num.MustU128FromLittleEndian(bs)
	tm := number.Rsh(17).And(num.U128From32(math.MaxUint32))
	ip := number.Rsh(49).Sub(tm)
	port := uint16(number.And(num.U128From64(0xFFFFFF)).Sub(tm.And(num.U128From32(0xFFFF))).AsUint64())
	bs = make([]byte, 16)
	ip.PutLittleEndian(bs)
	bs = bs[:4]
	str := ""
	for i, b := range bs {
		str += strconv.Itoa(int(b))
		if i != 3 {
			str += "."
		}
	}
	return fmt.Sprint(str, ":", port)
}
