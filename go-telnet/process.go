package main

import (
	"fmt"
	"strings"
)

// 输入“@close”退出当前连接会话。
// 输入“@shutdown”终止服务器运行。
func processTelnetCommand(str string, exitChan chan int) bool {
	if strings.HasPrefix(str, "@close") {
		fmt.Println("Session closed")
		// 告诉外部需要断开连接
		return false
	} else if strings.HasPrefix(str, "@shutdown") {
		fmt.Println("Server shutdown")
		exitChan <- 1
		return false
	}
	fmt.Println(str)
	return true
}
