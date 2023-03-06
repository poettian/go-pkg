package main

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"192.168.22.131:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
	}
	defer cli.Close()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Put(ctx, "sample_key", "sample_value")
	cancel()
	if err != nil {
		// handle error!
	}
	// use the response
}
