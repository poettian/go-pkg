package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

// 简单用法
func simpleSend()  {
	// 定义消息
	m := gomail.NewMessage()
	m.SetHeader("From", "poettian@163.com")
	m.SetHeader("To", "poettian@gmail.com")
	m.SetAddressHeader("Cc", "zhiwei.tian@eeoa.com", "志伟")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("/tmp/girl.jpeg")

	// 定义smtp拨号器
	d := gomail.NewDialer("smtp.163.com", 465, "poettian", "WBOWYFJNWICJAOXA")

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
	}
}

// Daemon用法
func daemonSend() {
	ch := make(chan *gomail.Message)

	go func() {
		// 定义smtp拨号器
		d := gomail.NewDialer("smtp.163.com", 465, "poettian", "WBOWYFJNWICJAOXA")

		var s gomail.SendCloser
		var err error
		open := false

		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}
				// 拨号连接
				if !open {
					if s, err = d.Dial(); err != nil {
						panic(err)
					}
					open = true
				}
				// 发送邮件
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if err := s.Close(); err != nil {
						panic(err)
					}
					open = false
				}
			}
		}
	}()

	close(ch)
}

func main()  {
	daemonSend()
}