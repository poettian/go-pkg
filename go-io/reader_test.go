package main

import "testing"

func TestReadLine1(t *testing.T) {
	str := ReadLine1("Hello world!\nThis is a small test case.")
	t.Log(str)
}
