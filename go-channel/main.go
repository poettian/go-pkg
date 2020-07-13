package main

import (
	"fmt"
	"runtime"
	"time"
)

func putData(intChan chan int) {
	for i := 1; i <= 800000; i++ {
		intChan <- i
	}
	close(intChan)
}

func primeNum(intChan chan int, primeChan chan int, cpuNum int, exitFlag *int) {
	for {
		flag := true
		num, ok := <-intChan
		if !ok {
			break
		}
		for i := 2; i < num; i++ {
			if num%i == 0 {
				flag = false
				break
			}
		}
		if flag {
			primeChan <- num
		}
	}
	// 如果所有的素数已经找到，则退出
	*exitFlag++
	if *exitFlag == cpuNum {
		close(primeChan)
	}
	fmt.Println("我已经跑完了")
}

func main() {
	cpuNum := runtime.NumCPU()
	starTime := time.Now().Unix()
	intChan := make(chan int, 1000)
	primeChan := make(chan int, 200)
	var exitFlag int = 0

	// 启一个协程负责生产数据
	go putData(intChan)

	// 启N个协程负责消费数据（查找素数并放入通道）
	for i := 0; i < cpuNum; i++ {
		go primeNum(intChan, primeChan, cpuNum, &exitFlag)
	}

	// 从通道中取出找到的素数并打印
	for {
		num, ok := <- primeChan
		if !ok {
			break
		}
		fmt.Printf("素数：%d\n", num)
	}

	endTime := time.Now().Unix()
	fmt.Println("执行时间", endTime-starTime)

	fmt.Println("main执行完了")
}
