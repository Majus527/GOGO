package Common

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"
)

type InputInfo struct {
	Hosts    string
	HostFile string
}

func Banner() {
	// 定义暗绿色系
	colors := []color.Attribute{
		color.FgGreen,   // 基础绿
		color.FgHiGreen, // 亮绿
	}

	lines := []string{
		"   ____  ____  ____  ____  ____",
		"  |  _ \\|  _ \\|  _ \\|  _ \\|  _ \\",
		"  | |_) | |_) | |_) | |_) | |_) |",
		"  |  __/|  __/|  __/|  __/|  __/",
		"  |_|   |_|   |_|   |_|   |_|  ",
	}

	// 获取最长行的长度
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	// 创建边框
	topBorder := "┌" + strings.Repeat("─", maxLength+2) + "┐"
	bottomBorder := "└" + strings.Repeat("─", maxLength+2) + "┘"

	// 打印banner
	fmt.Println(topBorder)

	for lineNum, line := range lines {
		fmt.Print("│ ")
		// 使用对应的颜色打印每个字符
		c := color.New(colors[lineNum%2])
		c.Print(line)
		// 补齐空格
		padding := maxLength - len(line)
		fmt.Printf("%s │\n", strings.Repeat(" ", padding))
	}

	fmt.Println(bottomBorder)

	// 打印版本信息
	c := color.New(colors[1])
	c.Printf("      GO-GO Version: %s\n\n", "1.0.0")
}

func (input *InputInfo) Flag() {
	// 目标配置
	flag.StringVar(&input.Hosts, "h", "", "输入目标ip或域名(,分隔)")
	flag.StringVar(&input.HostFile, "hf", "", "地址列表, -hf ip.txt")

	flag.Parse()
}

func (input *InputInfo) AnalyseUserInput() error {
	if input.HostFile == "" {
		return nil
	}

	// 验证文件是否存在
	if _, err := os.Stat(input.HostFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", input.HostFile)
	}

	// 验证是否有读取权限
	if _, err := os.Open(input.HostFile); os.IsPermission(err) {
		return fmt.Errorf("permission denied: %s", input.HostFile)
	}

	// 读取文件
	hosts, err := ReadFile(input.HostFile)
	if err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	input.Hosts = strings.Join(hosts, ",")
	return nil
}
