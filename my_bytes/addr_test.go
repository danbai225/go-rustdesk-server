package my_bytes

import "testing"

func TestAddr(t *testing.T) {
	addr := encodeAddr("192.168.16.32:21116")
	decodeAddr(addr)
}
