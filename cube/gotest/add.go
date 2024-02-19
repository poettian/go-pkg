package gotest

import "fmt"

// Add is a function that takes two integers and returns
// the sum of them
func Add(a, b int) int {
	fmt.Println(a, b)
	return a + b
}
