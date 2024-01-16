package config

import "fmt"

var std string

func init() {
	std = "hello"
}

func Echo() {
	fmt.Println(std)
}

func SetOutput(s string) {
	std = s
}
