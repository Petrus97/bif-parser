package bayesnet

import (
	"fmt"
	"reflect"
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
	IndexToAssignment(values, C.Numvalues)
}

func IndexToAssignment(values []int, cval []int) {
	fmt.Println("Before", values)
	newvct := make([][]int, 0)
	newvct = append(newvct, values)
	newvct = append(newvct, values)
	fmt.Println(newvct)
	fmt.Println(cval)
	vect := make([]int, 0)
	vect = append(vect, 1)
	for i := 0; i < len(cval)-1; i++ {
		k := vect[i] * cval[i]
		vect = append(vect, k)
	}
	fmt.Println(vect)
	numerator := make([][]int, len(cval))
	for l := range numerator {
		numerator[l] = make([]int, len(values))
	}
	for i, e := range vect {
		mult := newvct[i]
		for j := 0; j < len(mult); j++ {
			fmt.Println(mult[j], e, "=", mult[j]/e)
			numerator[i][j] = mult[j] / e
		}
	}
	fmt.Println(numerator)
	denominator := make([][]int, len(cval))
	for i := range denominator {
		denominator[i] = make([]int, len(values))
	}
	for i := 0; i < len(cval); i++ {
		denominator[i] = cval
	}
	fmt.Println(denominator)
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
