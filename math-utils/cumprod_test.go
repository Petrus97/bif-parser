package mathutils

import (
	"testing"
)

func TestCreateNewSlice(t *testing.T) {
	a := []int{0, 1, 2, 3, 4, 5, 6, 7}
	b := []int{0, 0, 0, 0, 0, 8, 9, 10}

	expcted := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	result := CreateNewSlice(a, b[5:])
	for i := 0; i < len(expcted); i++ {
		if expcted[i] != result[i] {
			t.Errorf("expected %d, got %d", expcted[i], result[i])
		}
	}
}

func TestCumprod(t *testing.T) {
	vect := NewVector([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12})
	result := Cumprod(vect.Data)
	expected := []int{1, 2, 6, 24, 120, 720, 5040, 40320, 362880, 3628800, 39916800, 479001600}
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Errorf("expected %d, got %d", expected[i], result[i])
		}
	}
}

func TestCumprodAndCreateSlice(t *testing.T) {
	vect := NewVector([]int{2, 2, 2, 2})
	result := Cumprod(CreateNewSlice([]int{1}, vect.Data[:len(vect.Data)-1]))
	expected := []int{1, 2, 4, 8}
	for i := 0; i < len(result); i++ {
		if result[i] != expected[i] {
			t.Errorf("expected %d, got %d", expected[i], result[i])
		}
	}
}
