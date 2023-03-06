package main

import (
	"fmt"
	"io"
)

func main() {
	var r io.Reader
	fmt.Println(r == nil)
}
