package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/Users/tianzhiwei/tmp/foo.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	log.Printf("Hello,%s", "world")
}
