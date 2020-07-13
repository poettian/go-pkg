package go_context

import (
	"context"
	"fmt"
	"time"
)

func monitor(ctx context.Context, number int)  {
	for {
		select {
		case <- ctx.Done():
			fmt.Printf("监控器%v，监控结束。\n", number)
			return
		default:
			// 获取 item 的值
			value := ctx.Value("item")
			fmt.Printf("监控器%v，正在监控 %v \n", number, value)
			time.Sleep(2 * time.Second)
		}
	}
}

func run() {
	ctx01, cancel1 := context.WithCancel(context.Background())
	ctx02, _ := context.WithTimeout(ctx01, 1* time.Second)
	ctx03 := context.WithValue(ctx02, "item", "CPU")

	defer cancel1()

	for i :=1 ; i <= 5; i++ {
		go monitor(ctx03, i)
	}

	time.Sleep(5  * time.Second)
	if ctx02.Err() != nil {
		fmt.Println("监控器取消的原因: ", ctx02.Err())
	}

	fmt.Println("主程序退出！！")
}