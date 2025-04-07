package Core

import (
	"fmt"
	"github.com/fatih/color"
	"gogo/Common"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 获取用户输入
func GetUserInput(input *Common.InputInfo) {
	input.Flag()
	err := input.AnalyseUserInput() // 解析用户的输入，把数据放入input中
	if err != nil {
		fmt.Println(err)
		return
	}
}

func ScanProcess() {
	// 获取用户输入并传入input中Common.InputInfo
	input := new(Common.InputInfo)
	GetUserInput(input)
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
	hostInfo := new(HostInfo)
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
	ps := NewPortScanner(hostInfo.HostAlive, ports, 3*time.Second)
	// ps := Core.NewPortScanner(hostInfo.HostAlive, []string{"445", "135", "3306", "7890", "10000", "7680", "12345"}, 1*time.Second)
	ps.TCPScan()
	fmt.Println("总共加载端口：", len(ports))
	color.Cyan("\n端口扫描结果如下：")
	for ip, ports := range ps.IpPortAlive {
		if len(ports) == 0 {
			continue
		}
		color.Green(ip + ": ")
		fmt.Println("总数：", len(ports))
		fmt.Println(ports)
	}
	color.Cyan("\n端口扫描结果汇总：")
	ipPort := []string{}
	for ip, ports := range ps.IpPortAlive {
		for _, port := range ports {
			fmt.Println(ip + ":" + port)
			ipPort = append(ipPort, ip+":"+port)
		}
	}
	// 将结果写入文件
	err := Common.WriteArrayToFile(Common.PortPath, ipPort)
	if err != nil {
		fmt.Println(err)
	}
	color.Cyan("\n已将结果写入文件" + Common.PortPath)

	// 指纹识别
	var wg sync.WaitGroup
	color.Cyan("\n进行指纹识别...")
	sem := make(chan struct{}, 999) // 控制最大并发数
	for ip, ports := range ps.IpPortAlive {
		for _, port := range ports {
			sem <- struct{}{} // 获取信号量
			wg.Add(1)         // 计数器+1
			go func(ip, port string) {
				defer func() {
					<-sem     // 释放信号量
					wg.Done() // 确保计数器-1
				}()
				res := MatchPortFinger(ip, port)
				if res != "Unknown" {
					ps.mu.Lock() // 加锁保护输出
					fmt.Println(ip + ":" + port + " is " + res)
					ps.mu.Unlock()
				}
			}(ip, port) // 传递当前 ip/port 的副本
		}
	}
	wg.Wait() // 等待所有 goroutine 完成

}
