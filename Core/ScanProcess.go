package Core

import (
	"fmt"
	"gogo/Common"
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
