package main

import "os"

func main()  {
	exitChan := make(chan int)

	go server("127.0.0.1:7001", exitChan)

	// 通道阻塞, 等待接收返回值
	code := <-exitChan

	// 标记程序返回值并退出
	os.Exit(code)
}
