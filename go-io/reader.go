package main

import (
	"bufio"
	"fmt"
	"strings"
)

func ReadLine1(s string) string {
	bufioReader := bufio.NewReader(strings.NewReader(s))
	bytes, err := bufioReader.ReadBytes('\n')
	fmt.Println(bytes, err)
	bytes, err = bufioReader.ReadBytes('\n')
	fmt.Println(bytes, err)
	return string(bytes)
}
