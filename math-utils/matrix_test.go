package mathutils

import "testing"

func TestTransposeVector(t *testing.T) {
	v := NewVector([]int{1, 2, 3, 4, 5, 6})
	transpose := v.T()

	if (v.M != transpose.N) || (v.N != transpose.M) {
		t.Error("Not transpose")
	}
}

func TestMatrixMod(t *testing.T) {
	numerator := NewMatrix(8, 3)
	numerator.Data = []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 0, 1, 1, 2, 2, 3, 3,
		0, 0, 0, 0, 1, 1, 1, 1,
	}
	denominator := NewMatrix(8, 3)
	denominator.Data = []int{
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
		2, 2, 2, 2, 2, 2, 2, 2,
	}
	result, err := MatrixMod(numerator, denominator)
	if err != nil {
		t.Error(err)
	}
	expected := NewMatrix(8, 3)
	expected.Data = []int{
		0, 1, 0, 1, 0, 1, 0, 1,
		0, 0, 1, 1, 0, 0, 1, 1,
		0, 0, 0, 0, 1, 1, 1, 1,
	}
	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != result.Data[i] {
			t.Errorf("Error, expected %d got %d", expected.Data[i], result.Data[i])
		}
	}
}

func TestMatrixDivision(t *testing.T) {
	repeatI := NewMatrix(8, 3)
	repeatI.Data = []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 1, 2, 3, 4, 5, 6, 7,
	}
	repeatCprodD := NewMatrix(8, 3)
	repeatCprodD.Data = []int{
		1, 1, 1, 1, 1, 1, 1, 1,
		2, 2, 2, 2, 2, 2, 2, 2,
		4, 4, 4, 4, 4, 4, 4, 4,
	}
	result, err := MatrixDivision(repeatI, repeatCprodD)
	if err != nil {
		t.Error(err)
	}
	expected := NewMatrix(8, 3)
	expected.Data = []int{
		0, 1, 2, 3, 4, 5, 6, 7,
		0, 0, 1, 1, 2, 2, 3, 3,
		0, 0, 0, 0, 1, 1, 1, 1,
	}
	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != result.Data[i] {
			t.Errorf("Error at %d, expected %d got %d", i, expected.Data[i], result.Data[i])
		}
	}
}

func TestMatrixDivision2(t *testing.T) {
	repeatI := NewMatrix(4, 2)
	repeatI.Data = []int{
		0, 1, 2, 3,
		0, 1, 2, 3,
	}
	repeatCprodD := NewMatrix(4, 2)
	repeatCprodD.Data = []int{
		1, 1, 1, 1,
		2, 2, 2, 2,
	}
	result, err := MatrixDivision(repeatI, repeatCprodD)
	if err != nil {
		t.Error(err)
	}
	expected := NewMatrix(4, 2)
	expected.Data = []int{
		0, 1, 2, 3,
		0, 0, 1, 1,
	}
	for i := 0; i < len(expected.Data); i++ {
		if expected.Data[i] != result.Data[i] {
			t.Errorf("Error at %d, expected %d got %d", i, expected.Data[i], result.Data[i])
		}
	}
}
