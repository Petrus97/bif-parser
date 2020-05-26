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
	if torepeat.M == 1 { // if is row vector
		idx := 0
		for i := 0; i < torepeat.N; i++ {
			tmp := make([]int, nrows)
			for j := 0; j < nrows; j++ {
				tmp[j] = torepeat.Data[i]
			}
			copy(m.Data[idx:], tmp)
			idx += len(tmp)
		}
	} else { // if is a coloumn vector
		idx := 0
		for idx <= dim {
			copy(m.Data[idx:], torepeat.Data)
			fmt.Println(m.Data)
			idx += len(torepeat.Data)
		}
	}
	return m
}
