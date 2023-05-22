package main

import (
	"context"
	"time"
)

func main() {
	background := context.Background()
	deadlineCtx, cancel := context.WithDeadline(background, time.Now().Add(10*time.Second))
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				println("cancel signal received")
				return
			default:
				println("hello")
				time.Sleep(1 * time.Second)
			}
		}
	}(deadlineCtx)
	<-time.After(5 * time.Second)
	cancel()
	<-time.After(2 * time.Second)
}
