package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {
		// 设置响应头
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// 写事件流开头
		fmt.Fprintf(w, "data: start\n\n")

		// 每隔1秒写入事件
		ticker := time.NewTicker(1 * time.Second)
		go func() {
			for range ticker.C {
				// 写事件
				fmt.Fprintf(w, "data: Message at %v\n\n", time.Now().Format(time.RFC3339))

				// 触发客户端重新连接
				fmt.Fprintf(w, ":reset\n\n")
			}
		}()
	})

	http.ListenAndServe(":8000", nil)
}
