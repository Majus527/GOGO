package Core

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type PortScanner struct {
	target      []string            // 目标地址
	ports       []string            // 要扫描的端口列表
	IpPortAlive map[string][]string // 存活列表
	timeout     time.Duration       // 超时时间
	wg          sync.WaitGroup
	mu          sync.Mutex // 添加互斥锁
}

// 新建一个PortScanner对象
func NewPortScanner(target []string, ports []string, timeout time.Duration) *PortScanner {
	return &PortScanner{
		target:      target,
		ports:       ports,
		IpPortAlive: make(map[string][]string),
		timeout:     timeout,
	}
}

// TCP扫描
func (ps *PortScanner) TCPScan() {
	sem := make(chan struct{}, 999) // 控制最大并发数

	for _, ip := range ps.target {
		ps.IpPortAlive[ip] = []string{}

		for _, port := range ps.ports {
			sem <- struct{}{} // 获取信号量
			ps.wg.Add(1)

			go func(ip string, port string) {
				defer func() {
					<-sem // 释放信号量
					ps.wg.Done()
				}()
				ps.TCPCheckPort(ip, port)
			}(ip, port)
		}
	}
	ps.wg.Wait()
	close(sem)
}

func (ps *PortScanner) TCPCheckPort(ip string, port string) {

	address := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.DialTimeout("tcp", address, ps.timeout)
	if err != nil {
		// fmt.Println(address, "is not open")
		return
	}
	ps.mu.Lock()
	ps.IpPortAlive[ip] = append(ps.IpPortAlive[ip], port)
	ps.mu.Unlock()
	fmt.Println(address, "is open")
	if conn != nil {
		conn.Close()
	}
}
