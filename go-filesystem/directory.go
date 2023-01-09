package go_filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetRuntimeDir 返回调用者所在文件目录路径，该函数返回的是编译时机器的路径而非执行时机器的路径
func GetRuntimeDir() string {
	// runtime.Caller 可以返回函数调用栈的某一层的程序计数器、文件信息、行号
	// 0 代表当前函数，也是调用runtime.Caller的函数。1 代表上一层调用者，以此类推。
	_, fileName, _, _ := runtime.Caller(1)
	pos := strings.LastIndex(fileName, "/")
	return fileName[:pos]
}

// GetExecutableDir1 返回程序所在目录路径
func GetExecutableDir1() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir
}

// GetExecutableDir2 返回程序所在目录路径，建议优先使用这种方式
func GetExecutableDir2() string {
	path, _ := os.Executable()
	return filepath.Dir(path)
}

// GetWorkDir 返回当前工作目录路径
func GetWorkDir() string {
	dir, _ := os.Getwd()
	return dir
}
