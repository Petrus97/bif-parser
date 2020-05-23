package bayesnet

import (
	"fmt"
	"reflect"

	mu "github.com/Petrus97/bif-parser/math-utils"
)

// Factor data type is the product of CPTs of nodes. ex: (X1)-->(X2)<--(X3) Φ(1, 2, 3) = P(X2 | X1, X3)*P(X1)*P(X3)
type Factor struct {
	Var       []*Node
	Numvalues []int
	Values    []float64
}

// CreateFactor takes in input a node of bayesian network and build a factor
func CreateFactor(n *Node) *Factor {
	factor := Factor{
		Var:       make([]*Node, 0),
		Numvalues: make([]int, 0),
		Values:    make([]float64, 0),
	}
	factor.Var = append(factor.Var, n)
	factor.Numvalues = append(factor.Numvalues, n.Numvalues)
	for _, parent := range n.Parents {
		factor.Var = append(factor.Var, parent)
		factor.Numvalues = append(factor.Numvalues, parent.Numvalues)
	}
	factor.Values = append(factor.Values, n.CPT...)
	return &factor
}

func FactorProduct(A *Factor, B *Factor) {
	intersection := intersect(A, B)
	if len(intersection) < 0 {
		return
	}
	C := Factor{}
	C.Var = union(A, B, intersection)
	C.Numvalues = make([]int, len(C.Var))
	mapA := mapIndexValue(A, &C)
	mapB := mapIndexValue(B, &C)
	for i, index := range mapA {
		C.Numvalues[index] = A.Numvalues[i]
	}
	for i, index := range mapB {
		C.Numvalues[index] = B.Numvalues[i]
	}
	var nval int = 1
	for _, i := range C.Numvalues {
		nval *= i
	}
	C.Values = make([]float64, nval)
	fmt.Println(mapA, mapB, C)
	values := make([]int, 0)
	for i := 0; i < len(C.Values); i++ {
		values = append(values, i)
	}
	assignments := IndexToAssignment(values, C.Numvalues)
	fmt.Println(assignments)

}

func IndexToAssignment(factorvalues []int, factornval []int) *mu.Matrix {
	valuevector := mu.NewVector(factorvalues)
	nvalvector := mu.NewVector(factornval)
	repeatI := mu.Repmat(valuevector.T(), 1, len(factornval))
	fmt.Println("repeat_I", repeatI)
	cprod := mu.Cumprod(mu.CreateNewSlice([]int{1}, nvalvector.Data[:len(nvalvector.Data)-1]))
	cprodvector := mu.NewVector(cprod)
	repeatD := mu.Repmat(cprodvector, len(factorvalues), 1)
	fmt.Println("repeat_D", repeatD)
	numerator, err := mu.MatrixDivision(repeatI, repeatD)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("numerator", numerator)
	denominator := mu.Repmat(nvalvector, len(factorvalues), 1)
	fmt.Println("den", denominator)
	indexes, err := mu.MatrixMod(numerator, denominator)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("indexes", indexes)
	return indexes
	// fmt.Println("Before", values)
	// newvct := make([][]int, 0)
	// myvector := make([][]int, len(values))
	// for i := 0; i < len(values); i++ {
	// 	myvector[i] = make([]int, 0)
	// }
	// for i := 0; i < len(values); i++ {
	// 	if i%len(cval) == 0 {
	// 		myvector[i] = values[:len(cval)]
	// 	} else {
	// 		myvector[i] = values[len(cval):]
	// 	}

	// }
	// fmt.Println("myvector", myvector)
	// newvct = append(newvct, values)
	// newvct = append(newvct, values)
	// fmt.Println("newvct", newvct)
	// fmt.Println("cval", cval)
	// vect := make([]int, 0)
	// vect = append(vect, 1)
	// for i := 0; i < len(cval)-1; i++ {
	// 	k := vect[i] * cval[i]
	// 	vect = append(vect, k)
	// }
	// fmt.Println("vect", vect)
	// numerator := make([][]int, len(cval))
	// for l := range numerator {
	// 	numerator[l] = make([]int, len(values))
	// }
	// newnumerator := make([][]int, len(values))
	// for i := range newnumerator {
	// 	newnumerator[i] = make([]int, len(cval))
	// }

	// // [[0 0] [1 0] [2 1] [3 1]]
	// /*
	// 	0 0 = 0 / 1  0 / 2
	// 	1 0 = 1 / 1  1 / 2
	// 	2 1 = 2 / 1  2 / 2
	// 	3 1 = 3 / 1  3 / 2
	// */
	// fmt.Println("newnumerator", newnumerator)
	// for i, e := range vect {
	// 	mult := newvct[i]
	// 	for j := 0; j < len(mult); j++ {
	// 		fmt.Println(mult[j], e, "=", mult[j]/e)
	// 		numerator[i][j] = mult[j] / e
	// 	}
	// }
	// fmt.Println("numerator", numerator)
	// denominator := make([][]int, len(cval))
	// for i := range denominator {
	// 	denominator[i] = make([]int, len(values))
	// }
	// for i := 0; i < len(cval); i++ {
	// 	denominator[i] = cval
	// }
	// fmt.Println("denominator", denominator)
}

func AssignmentToIndex(assignments *mu.Matrix, factorcard []int) {
	//cardvect := mu.NewVector(factorcard)
	//cprodD := mu.Cumprod(1, cardvect)

}

func intersect(A *Factor, B *Factor) []*Node {
	intersection := make([]*Node, 0)
	for _, node := range A.Var {
		if containsNode(B, node) {
			intersection = append(intersection, node)
		}
	}
	return intersection
}

func containsNode(f *Factor, n *Node) bool {
	for _, node := range f.Var {
		if node == n {
			return true
		}
	}
	return false
}

func union(A *Factor, B *Factor, intersection []*Node) []*Node {
	union := make([]*Node, 0)
	union = append(union, A.Var...)
	union = append(union, B.Var...)

	for _, node := range intersection {
		if ok, index := contains(union, node); ok == true {
			union = removeIndex(union, index)
		}
	}
	return union
}

func removeIndex(union []*Node, index int) []*Node {
	return append(union[:index], union[index+1:]...)
}

func contains(a interface{}, e interface{}) (bool, int) {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true, i
		}
	}
	return false, -1
}

func mapIndexValue(f *Factor, C *Factor) []int {
	pos := make([]int, 0)
	for _, variable := range f.Var {
		fmt.Println(variable)
		for index, cvar := range C.Var {
			if cvar == variable {
				pos = append(pos, index)
			}
		}
	}
	return pos
}
