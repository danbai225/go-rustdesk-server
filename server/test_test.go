package server

import (
	logs "github.com/danbai225/go-logs"
	"testing"
	"time"
)

func TestRingMsgMsg(t *testing.T) {
	go func() {
		time.Sleep(time.Second * 5)
		r.Put(&ringMsg{
			ID:      "1",
			Type:    "2",
			TimeOut: 2,
			InsTime: time.Now(),
			Val:     123,
		})
	}()
	form := getMsgForm("1", "2", 3)
	logs.Info(form)
}
