package go_thread

import (
	"log"
	"time"
)

// 这个例子用来看内部的goroutine是否会随着caller的结束而结束
func CheckInnerGoroutineExit() {
	runInnerGoroutine()
	log.Println("app run out of checkIfGoroutineExit()")
	time.Sleep(time.Second * 30)
}

func runInnerGoroutine() {
	// 当前函数退出不影响内部定义的goroutine的执行
	// 结合《go并发编程实战》可知，G会被放到P的队列中等待运行
	go func() {
		ticker := time.NewTicker(time.Second)
		after := time.After(time.Second * 5)
		for {
			select {
			case t := <-ticker.C:
				log.Println(t.Format("2006-01-02 15:04:05"))
			}
			// 循环到第二次这里就阻塞住了
			v, ok := <-after
			log.Println("time over", v, ok)
		}
	}()
	log.Println("start time ticker")
}
