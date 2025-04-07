package Core

import (
	_ "embed"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
	"time"
)

//go:embed nmap-service-probes.txt
var ProbeString string

func ParseNmapProbeString(input string) ([]string, error) {
	// 正则表达式匹配三部分：
	// 1. 服务名称 (mysql)
	// 2. 匹配模式 (m|...| 之间的内容)
	// 3. 其他属性 (s/... i/... cpe:... 等)
	re := regexp.MustCompile(`^match\s+(\w+)\s+m\|([^|]+)\|(.*)$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid input format")
	}

	return []string{matches[1], matches[2], matches[3]}, nil
}

func TCPCheckPort(ip string, port string) (net.Conn, error) {

	address := fmt.Sprintf("%s:%s", ip, port)
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		// fmt.Println(address, "is not open")
		return nil, err
	}
	// fmt.Println(address, "is open")

	return conn, nil
}

func MatchPortFinger(ip string, port string) (FingerName string) {
	fmt.Println("正在识别：", ip, ":", port)
	FingerName = "Unknown"
	// 按换行符分割成行
	lines := strings.Split(ProbeString, "\n")
	conn, err := TCPCheckPort(ip, port)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	data, err := io.ReadAll(conn)
	// fmt.Println("data:", string(data))

	// 遍历每一行，过滤空白行
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line) // 去除首尾空白字符
		if trimmedLine != "" {                 // 检查是否非空
			// fmt.Println(trimmedLine)
			res, err := ParseNmapProbeString(trimmedLine)
			if err != nil {
				fmt.Println(err)
				return
			}
			pattern := res[1]
			// fmt.Println("pattern:", pattern)
			// 编译正则表达式
			regex, err := regexp.Compile(pattern)
			if err != nil {
				fmt.Println("正则表达式编译失败:", err)
				return
			}

			if regex.MatchString(string(data)) {
				FingerName = res[0]
				// fmt.Println("匹配成功！")
				return
			} else {
				continue
			}

		}
	}
	return
}
