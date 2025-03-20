package main

import (
	"fmt"
	"gogo/Common"
	"gogo/Core"
)

func main() {
	Common.Banner()
	// 获取用户输入
	input := new(Common.InputInfo)
	input.Flag()
	err := input.AnalyseUserInput()
	if err != nil {
		fmt.Println(err)
		return
	}
	hosts := Common.InputHostsHandler(input.Hosts)

	// 初始化
	fmt.Println("正在进行去重操作...")
	hosts = Common.RemoveRepFromMap(hosts) // 去重
	fmt.Println("域名和ip解析中...")
	ips, domains := Common.SeparateIPsAndDomains(hosts)

	fmt.Println("域名如下：")
	fmt.Println(domains)

	fmt.Println("ip如下：")
	fmt.Println(ips)

	fmt.Println("进行探活操作...")

	hostInfo := new(Core.HostInfo)
	hostInfo.InitHostInfo(hosts, domains, ips)

	// 探测存活
	hostInfo.PingHosts()

	fmt.Println("存活资产如下：")
	fmt.Println(hostInfo.HostAlive)

}
