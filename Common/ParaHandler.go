package Common

import (
	"net"
	"strings"
)

// 字符串数组去重
func RemoveRepFromMap(slc []string) []string {
	result := []string{}         //存放返回的不重复切片
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// isIP checks if a given string is a valid IP address (both IPv4 and IPv6)
func isIP(s string) bool {
	return net.ParseIP(s) != nil
}

// 区分IP和域名
func SeparateIPsAndDomains(input []string) ([]string, []string) {
	var ips []string
	var domains []string

	for _, str := range input {
		str = strings.TrimSpace(str) // Remove any leading/trailing whitespace
		if isIP(str) {
			ips = append(ips, str)
		} else {
			domains = append(domains, str)
		}
	}

	return ips, domains
}

func InputHostsHandler(input string) (hosts []string) {
	hosts = strings.Split(input, ",")
	return
}
