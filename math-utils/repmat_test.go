package mathutils

import "testing"

func TestRepmat(t *testing.T) {
	slicetest := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	vtest := NewVector(slicetest)
	result := Repmat(vtest.T(), 1, 4)

	expected := new(Matrix)
	expected.M = 16
	expected.N = 4
	expected.Data = []int{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16,
	}
	if expected.M != result.M {
		t.Errorf("Expected nrows %d, got %d", expected.M, result.M)
	}
	if expected.N != result.N {
		t.Errorf("Expected ncols %d, got %d", expected.N, result.N)
	}
	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != result.Data[i] {
			t.Errorf("Expected data at %d: %d, got %d", i, expected.N, result.N)
		}
	}
}
