package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	// max execute time: 5 seconds
	if err := slowOption(ctx); err != nil {
		fmt.Println(err)
	}
}

func slowOption(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("time over")
			return ctx.Err()
		case <-time.After(time.Second):
			fmt.Println("程序运行中...")
		}
	}
}
