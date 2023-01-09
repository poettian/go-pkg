package go_filesystem

import "testing"

func TestGetExecutableDir1(t *testing.T) {
	dir := GetExecutableDir1()
	t.Log(dir)
}

func TestGetExecutableDir2(t *testing.T) {
	path := GetExecutableDir2()
	t.Log(path)
}
