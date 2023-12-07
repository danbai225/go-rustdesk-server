package common

import (
	"fmt"
	logs "github.com/danbai225/go-logs"
	"net"
	"os"
	"reflect"
	"strings"
)

func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {
	return !IsDir(path)
}
func InSameSubnet(ip1Str, ip2Str, maskStr string) bool {
	ip1 := net.ParseIP(ip1Str)
	ip2 := net.ParseIP(ip2Str)
	mask := net.IPMask(net.ParseIP(maskStr).To4())

	if ip1 == nil || ip2 == nil {
		return false
	}

	network1 := ip1.Mask(mask)
	network2 := ip2.Mask(mask)
	return network1.Equal(network2)
}

func InSubnet(ip string) bool {
	if !(strings.HasPrefix(ip, "192.168.") || strings.HasPrefix(ip, "10.") || strings.HasPrefix(ip, "172.")) {
		return false
	}
	// 获取所有的网络接口
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	// 遍历所有的网络接口
	for _, i := range interfaces {
		// 获取接口的所有地址
		addrs, _ := i.Addrs()
		// 遍历接口的所有地址
		for _, addr := range addrs {
			// 检查地址是否是IPNet类型，IPNet类型的地址包含子网掩码信息
			if ipNet, ok := addr.(*net.IPNet); ok {
				// 获取本机的IP地址
				localIP := ipNet.IP.String()
				// 如果这个地址是你的内网IP地址
				if strings.HasPrefix(localIP, "192.168.") || strings.HasPrefix(localIP, "10.") || strings.HasPrefix(localIP, "172.") {
					// 打印出子网掩码
					logs.Debug("Subnet mask for IP", localIP, "is", ipNet.Mask.String(), "ip:", ip)
					// 检查ip地址是否在子网内
					if InSameSubnet(localIP, ip, ipNet.Mask.String()) {
						return true
					}
				}
			}
		}
	}
	return false
}
