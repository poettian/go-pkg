package gotest

func MakeSliceWithoutAlloc() []int {
	var newSlice []int

	for i := 0; i < 100000; i++ {
		newSlice = append(newSlice, i)
	}

	return newSlice
}

func MakeSliceWithPreAlloc() []int {
	newSlice := make([]int, 100000)

	for i := 0; i < 100000; i++ {
		newSlice[i] = i
	}

	return newSlice
}
