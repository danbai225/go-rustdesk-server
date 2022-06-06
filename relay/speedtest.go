package relay

import (
	logs "github.com/danbai225/go-logs"
	p "github.com/go-ping/ping"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/net"
	sp "github.com/showwin/speedtest-go/speedtest"
	"go-rustdesk-server/common"
	"strings"
	"time"
)

func testSpeed() (Download, Upload, pingR uint) {
	user, _ := sp.FetchUserInfo()
	serverList, _ := sp.FetchServers(user)
	targets, _ := serverList.FindServer([]int{})
	server := strings.Split(common.Conf.RegServer, ":")
	for _, s := range targets {
		_ = s.DownloadTest(false)
		_ = s.UploadTest(false)
		return uint(s.DLSpeed), uint(s.ULSpeed), ping(server[0])
	}
	return
}
func ping(host string) uint {
	pinger := p.New(host)
	pinger.SetPrivileged(false)
	pinger.Count = 3
	pinger.Timeout = time.Second
	err := pinger.Run()
	if err != nil {
		logs.Err(err)
	}
	statistics := pinger.Statistics()
	return uint(statistics.AvgRtt.Milliseconds())
}
func cpuTest() uint {
	percent, _ := cpu.Percent(time.Second, false)
	return uint(percent[0])
}
func netFlow() float64 {
	info, _ := net.IOCounters(true)
	max := net.IOCountersStat{}
	for _, v := range info {
		if v.BytesSent > max.BytesSent {
			max = v
		}
	}
	time.Sleep(time.Second)
	info, _ = net.IOCounters(true)
	for _, stat := range info {
		if stat.Name == max.Name {
			m := float64(stat.BytesSent - max.BytesSent)
			m *= 8
			m /= 1024 * 1024
			return m
		}
	}
	return 0
}
