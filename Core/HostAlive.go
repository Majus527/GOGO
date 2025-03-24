package Core

import (
	"bytes"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

type HostInfo struct {
	HostAll      []string
	HostAlive    []string
	IPsAll       []string
	DomainsAll   []string
	IPsAlive     []string
	DomainsAlive []string
	mu           sync.Mutex // 添加互斥锁
}

// 初始化HostAlive
func (h *HostInfo) InitHostInfo(hostAll []string, domainsAll []string, ipsAll []string) {
	h.HostAll = hostAll
	h.DomainsAll = domainsAll
	h.IPsAll = ipsAll
	h.HostAlive = make([]string, 0)
}

func (h *HostInfo) PingHosts() {
	var wg sync.WaitGroup
	for _, host := range h.HostAll {
		wg.Add(1) // 为每个 ping 任务增加一个 WaitGroup 计数
		go func(host string) {
			defer wg.Done() // 任务完成后减少一个 WaitGroup 计数
			if pingHost(host) {
				h.mu.Lock() // 在修改切片前加锁
				h.HostAlive = append(h.HostAlive, host)
				h.mu.Unlock() // 在修改切片后解锁
			}
		}(host)
	}
	wg.Wait() // 等待所有 goroutine 完成
}

// 使用系统的ping命令检测主机存活状态
func pingHost(host string) bool {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "ping -n 2 -w 10 "+host+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 10 "+host+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	default: //linux
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 10 "+host+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	}
	outinfo := bytes.Buffer{}
	command.Stdout = &outinfo
	err := command.Start()
	if err != nil {
		return false
	}
	if err = command.Wait(); err != nil {
		return false
	}
	if strings.Contains(outinfo.String(), "TTL expired in transit") {
		return false
	}
	if strings.Contains(outinfo.String(), "true") {
		return true
	} else {
		return false
	}

}
