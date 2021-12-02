package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func readFrom(reader io.Reader, num int) ([]byte, error) {
	buffer := make([]byte, num)
	n, err := reader.Read(buffer)

	return buffer[:n], err
}

func main()  {
	data, err := readFrom(os.Stdin, 11)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(data))
}
