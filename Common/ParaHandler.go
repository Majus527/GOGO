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

// 将用户输入的host全部返回
func InputHostsHandler(input *InputInfo) (hosts []string) {
	if input.Hosts != "" {
		hosts = strings.Split(input.Hosts, ",")
	} else if len(input.HostFile) != 0 {
		hosts = append(hosts, input.HostFile...)
	}
	return
}

// 将用户输入的port进行格式化，返回一个数组
func InputPortsHandler(input *InputInfo) (ports []string) {
	if input.Ports != "" {
		ports = strings.Split(input.Ports, ",")
	} else {
		ports = []string{"21", "22", "23", "25", "53", "80", "110", "135", "139", "143", "443", "445", "1433", "1521", "3306", "3389", "5432", "5900", "8080", "8081", "9200", "9300", "11211", "27017", "6379", "7001", "8000", "8009", "8088", "9000", "9042", "9200", "9300", "10000", "27017"}
	}
	return
}
