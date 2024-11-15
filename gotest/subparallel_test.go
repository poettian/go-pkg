package gotest

import (
	"testing"
	"time"
)

func parallelTest1(t *testing.T) {
	t.Parallel()
	time.Sleep(3 * time.Second)
	// do some testing
}

func parallelTest2(t *testing.T) {
	t.Parallel()
	time.Sleep(2 * time.Second)
	// do some testing
}

func parallelTest3(t *testing.T) {
	t.Parallel()
	time.Sleep(1 * time.Second)
	// do some testing
}

func TestSubParallel(t *testing.T) {
	// setup
	t.Logf("setup")

	// t.Run会启动一个goroutine来执行func，在这之前会新生成一个子t并传到func中
	t.Run("group", func(t *testing.T) {
		t.Run("Test1", parallelTest1)
		t.Run("Test2", parallelTest2)
		t.Run("Test3", parallelTest3)
	})

	// teardown
	t.Logf("teardown")
}

//func TestMain(m *testing.M) {
//	println("TestMain setup.")
//
//	retCode := m.Run() // 执行测试，包括单元测试、性能测试和示例测试
//
//	println("TestMain teardown.")
//
//	os.Exit(retCode)
//}
