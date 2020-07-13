package main

import (
	"fmt"
	"net"
)

// 监听器
func server(address string, exitChan chan int) {
	// 根据给定地址进行监听
	l, err := net.Listen("tcp", address)

	if err != nil {
		fmt.Println(err)
		exitChan <- 1
	}

	fmt.Println("listen:" + address)

	defer l.Close()

	for {
		conn, err := l.Accept()

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		go handleSession(conn, exitChan)
	}
}
