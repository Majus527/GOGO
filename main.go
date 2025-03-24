package main

import (
	"fmt"
	"github.com/fatih/color"
	"gogo/Common"
	"gogo/Core"
	"strconv"
	"strings"
	"time"
)

func main() {
	Common.Banner()

	// 获取用户输入并传入input中Common.InputInfo
	input := new(Common.InputInfo)
	Core.GetUserInput(input)
	// 拿到用户输入的所有hosts
	hosts := Common.InputHostsHandler(input)

	if len(hosts) == 0 || len(hosts) == 1 && hosts[0] == "" { // 可能会出现{"":""}的情况
		color.Red("做啥？？？？")
		return
	}

	// 初始化
	fmt.Println("正在进行去重操作...")
	hosts = Common.RemoveRepFromMap(hosts) // 去重
	fmt.Println("域名和ip解析中...")
	ips, domains := Common.SeparateIPsAndDomains(hosts)

	color.Cyan("\n域名如下：")
	fmt.Println("总数：", len(domains))
	fmt.Println(domains)
	color.Cyan("\nip如下：")
	fmt.Println("总数：", len(ips))
	fmt.Println(ips)

	fmt.Println("进行探活操作...")

	// 引入hostInfo变量
	hostInfo := new(Core.HostInfo)
	hostInfo.InitHostInfo(hosts, domains, ips)

	// 探测存活
	hostInfo.PingHosts()

	color.Cyan("\n存活资产如下：")
	fmt.Println("总数：", len(hostInfo.HostAlive))
	fmt.Println(hostInfo.HostAlive)

	// 进行端口扫描
	// 处理端口参数
	fmt.Println("进行端口扫描操作...")
	ports := []string{}
	switch input.Ports {
	case "1": // 只扫主要web端口
		ports = strings.Split(Common.WebPorts, ",")
	case "2": // 精简主要端口
		ports = strings.Split(Common.MainPorts, ",")
	case "3": // 全端口
		ports = make([]string, 65535)
		for i := range ports {
			ports[i] = strconv.Itoa(i + 1)
		}
	default: // 默认扫web
		ports = strings.Split(Common.WebPorts, ",")
	}
	ps := Core.NewPortScanner(hostInfo.HostAlive, ports, 1*time.Second)
	// ps := Core.NewPortScanner(hostInfo.HostAlive, []string{"445", "135", "3306", "7890", "10000", "7680", "12345"}, 1*time.Second)
	ps.TCPScan()
	fmt.Println("总共加载端口：", len(ports))
	color.Cyan("\n端口扫描结果如下：")
	for ip, ports := range ps.IpPortAlive {
		color.Green(ip + ": ")
		fmt.Println("总数：", len(ports))
		fmt.Println(ports)
	}
	color.Cyan("\n端口扫描结果汇总：")
	for ip, ports := range ps.IpPortAlive {
		for _, port := range ports {
			fmt.Println(ip + ":" + port)
		}
	}
}
