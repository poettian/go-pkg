package main

import (
	"os"
)

func OpenFile(path string) (*os.File, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return fd, nil
}
