package my_bytes

import "testing"

func TestAddr(t *testing.T) {
	str := "192.168.16.32:21116"
	addr := EncodeAddr(str)
	s := DecodeAddr(addr)
	if s != str {
		t.Error()
	}
}
