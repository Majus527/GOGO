package Common

import (
	"bufio"
	"os"
)

// 从文件中读取数据，并返回一个字符串切片
func ReadFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var hosts []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hosts = append(hosts, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return hosts, nil
}

// WriteArrayToFile 将字符串数组逐行写入文件
func WriteArrayToFile(filename string, data []string) error {
	// 打开文件，如果不存在则创建
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close() // 确保文件最终关闭

	// 逐行写入数组数据
	for _, line := range data {
		_, err := file.WriteString(line + "\n") // 写入每一行并添加换行符
		if err != nil {
			return err
		}
	}

	return nil
}
