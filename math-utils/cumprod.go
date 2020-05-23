package mathutils

import "fmt"

// Cumprod multiply each element of a slice with his next
func Cumprod(slice []int) []int {
	res := make([]int, len(slice))
	i := 0
	product := 1
	for i < len(res) {
		res[i] = slice[i] * product
		product = res[i]
		i++
	}
	return res
}

// CreateSubSlice create a slice by
func CreateNewSlice(start []int, end []int) []int {
	fmt.Println("start", start, "end", end)
	newslice := make([]int, len(start)+len(end))
	copy(newslice[:len(start)], start[:])
	copy(newslice[len(start):], end[:])
	return newslice
}
