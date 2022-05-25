package common

var ipMap map[string]struct{}

func loadList() {
	ipMap = make(map[string]struct{})
	for _, s := range Conf.IpList {
		ipMap[s] = struct{}{}
	}
}
func InList(ip string) bool {
	if _, ok := ipMap[ip]; ok {
		return true
	}
	return false
}
func GetList() []string {
	ips := make([]string, len(ipMap))
	for k, _ := range ipMap {
		ips = append(ips, k)
	}
	return ips
}
func UpDataList(w bool, list []string) {
	Conf.WhiteList = w
	Conf.IpList = list
	Save()
	loadList()
}
