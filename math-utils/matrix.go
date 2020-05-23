package mathutils

import "fmt"

type Matrix struct {
	M, N int
	Data []int
}

type Vector struct {
	M, N int // vector dimensions
	Data []int
}

func NewMatrix(m, n int) (matrix *Matrix) {
	matrix = new(Matrix)
	matrix.M = m
	matrix.N = n
	matrix.Data = make([]int, m*n)
	return matrix
}

func MatrixDivision(numerator *Matrix, denominator *Matrix) (m *Matrix, err error) {
	if numerator.N != denominator.N {
		return nil, fmt.Errorf("dimension mismatch")
	}
	m = new(Matrix)
	m.M = numerator.M
	m.N = denominator.N
	dim := m.N * m.M
	m.Data = make([]int, dim)
	k := 0
	for i := 0; i < m.N; i++ {
		for j := k; j < k+m.M; j++ {
			// fmt.Println(numerator.Data[j], "/", denominator.Data[j], "=", numerator.Data[j]/denominator.Data[j])
			m.Data[j] = numerator.Data[j] / denominator.Data[j]
		}
		k += (m.M)
	}
	// for i := 0; i < dim; i++ {
	// 	fmt.Println(numerator.Data[i], "/", denominator.Data[i], "=", numerator.Data[i]/denominator.Data[i])
	// 	m.Data[i] = numerator.Data[i] / denominator.Data[i]
	// }
	return m, nil
}

func MatrixMod(numerator *Matrix, denominator *Matrix) (m *Matrix, err error) {
	if numerator.N != denominator.N {
		return nil, fmt.Errorf("dimension mismatch")
	}
	m = new(Matrix)
	m.M = numerator.M
	m.N = denominator.N
	dim := m.N * m.M
	m.Data = make([]int, dim)
	// k := 0
	// fmt.Println(m.M, m.N)
	// for i := 0; i < m.N; i++ { // on rows
	// 	for j := 0; j < m.M; j++ { // on cols
	// 		fmt.Println(numerator.Data[j], denominator.Data[i])
	// 		m.Data[k+j] = numerator.Data[j] % denominator.Data[i]
	// 	}
	// 	k += (m.M)
	// }
	for i := 0; i < dim; i++ {
		fmt.Println(numerator.Data[i], "/", denominator.Data[i], "=", numerator.Data[i]/denominator.Data[i])
		m.Data[i] = numerator.Data[i] % denominator.Data[i]
	}
	return m, nil
}

// NewVector builds a row vector
func NewVector(slice []int) (v *Vector) {
	v = new(Vector)
	v.Data = slice
	v.M = 1
	v.N = len(slice)
	return v
}

// T transpose a vector
func (v *Vector) T() (t *Vector) {
	t = new(Vector)
	t.Data = v.Data
	t.M = v.N
	t.N = v.M
	return t
}
