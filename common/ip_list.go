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
