package gotest_test

import (
	gotest "learn-packages/go-test"
	"testing"
)

func BenchmarkMakeSliceWithoutAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotest.MakeSliceWithoutAlloc()
	}
}

func BenchmarkMakeSliceWithPreAlloc(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gotest.MakeSliceWithPreAlloc()
	}
}
