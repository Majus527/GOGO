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
	for _, ip := range ps.target {
		ps.IpPortAlive[ip] = []string{}
		for _, port := range ps.ports {
			ps.wg.Add(1)
			go ps.TCPCheckPort(ip, port)
		}
	}
	ps.wg.Wait()
}

func (ps *PortScanner) TCPCheckPort(ip string, port string) {
	defer ps.wg.Done()

	address := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.DialTimeout("tcp", address, ps.timeout)
	if err != nil {
		return
	}
	ps.mu.Lock()
	ps.IpPortAlive[ip] = append(ps.IpPortAlive[ip], port)
	ps.mu.Unlock()
	if conn != nil {
		conn.Close()
	}
}
