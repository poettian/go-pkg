package main

import (
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/proc"
)

func main() {
	proc.AddShutdownListener(func() {
		println("shutdown")
	})
	proc.AddShutdownListener(func() {
		println("shutdown2")
	})
	time.Sleep(20 * time.Second)
	fmt.Println("start")
}
