package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	input := "date"
	hmd5 := md5.New()
	hmd5.Write([]byte(input))
	fmt.Printf("%x\n", hmd5.Sum(nil))
	input1 := "2023"
	hmd5.Write([]byte(input1))
	fmt.Printf("%x\n", hmd5.Sum(nil))
}
