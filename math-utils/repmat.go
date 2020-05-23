package mathutils

import "fmt"

func Repmat(torepeat *Vector, nrows, ncols int) (m *Matrix) {
	// fmt.Println("DEBUG")
	// fmt.Println(torepeat.Data)
	// fmt.Println("M", torepeat.M)
	// fmt.Println("N", torepeat.N)
	m = NewMatrix(torepeat.M*nrows, torepeat.N*ncols)
	dim := (torepeat.M * nrows) * (torepeat.N * ncols)
	fmt.Println(nrows, ncols)
	fmt.Println(m.M, m.N)
	fmt.Println(dim)
	idx := 0
	for idx <= dim {
		copy(m.Data[idx:], torepeat.Data)
		fmt.Println(m.Data)
		idx += len(torepeat.Data)
	}
	return m
}
