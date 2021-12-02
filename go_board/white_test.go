package go_board

import (
	"testing"
)

func TestStringToByte(t *testing.T) {
	StringToBytes([]byte("Hello World!"))
}

func TestCallMeTimer(t *testing.T) {
	CallMeTimer(5)
}

func TestBufioScanner(t *testing.T) {
	BufioScanner()
}

func TestGetRandomToken(t *testing.T) {
	token := GetRandomToken(32)
	t.Log(token)
}

func TestBoard(t *testing.T) {
	GetCpuNums()
}

