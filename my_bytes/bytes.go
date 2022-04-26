package my_bytes

import (
	"encoding/binary"
	"errors"
	"math"
)

func Decode(src []byte) (data []byte, err error) {
	if src == nil {
		return
	}
	headLen := uint((src[0] & 0x3) + 1)
	if uint(len(src)) < headLen {
		err = errors.New("dataLen<headLen")
		return
	}
	n := uint(src[0])
	if headLen > 1 {
		n |= uint(src[1]) << 8
	}
	if headLen > 2 {
		n |= uint(src[2]) << 16
	}
	if headLen > 3 {
		n |= uint(src[3]) << 24
	}
	n >>= 2
	if n > math.MaxUint {
		err = errors.New("too big packet")
		return
	}
	src = append(make([]byte, n), src...)
	data = src[(len(src) - int(n)):]
	return
}

func Encoder(src []byte) (data []byte, err error) {
	if src == nil {
		return
	}
	l := len(src)
	if l <= 0x3F {
		src = append([]byte{byte(l << 2)}, src...)
	} else if l <= 0x3FFF {
		temp := make([]byte, 2)
		binary.LittleEndian.PutUint16(temp, uint16(l<<2)|0x1)
		src = append(temp, src...)
	} else if l <= 0x3FFFFF {
		h := uint32(l<<2) | 0x2
		temp := make([]byte, 2)
		binary.LittleEndian.PutUint16(temp, uint16(h&0xFFFF))
		src = append(temp, src...)
		src = append([]byte{byte(h >> 16)}, src...)
	} else if l <= 0x3FFFFFFF {
		temp := make([]byte, 4)
		binary.LittleEndian.PutUint32(temp, uint32((l<<2)|0x3))
		src = append(temp, src...)
	} else {
		err = errors.New("overflow")
	}
	return src, nil
}
