package bayesnet

import (
	"fmt"
	"reflect"

	mu "github.com/Petrus97/bif-parser/math-utils"
	"github.com/fatih/color"
	"github.com/jinzhu/copier"
	ts "gorgonia.org/tensor"
)

// Factor data type is the product of CPTs of nodes. ex: (X1)-->(X2)<--(X3) Φ(1, 2, 3) = P(X2 | X1, X3)*P(X1)*P(X3)
// This is for a vectorized version. deprecated.
type Factor struct {
	Var       []*Node
	Numvalues []int
	Values    []float64
}

// FactorV2 data type is the product of CPTs of nodes. ex: (X1)-->(X2)<--(X3) Φ(1, 2, 3) = P(X2 | X1, X3)*P(X1)*P(X3)
type FactorV2 struct {
	Scope   []*Node
	Card    []int
	CPT     []float64
	Strides map[*Node]int
}

// CreateFactor takes in input a node of bayesian network and build a factor
// This is for a vectorized version. deprecated.
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

// CreateFactorV2 takes in input a node of bayesian network and build a factor
func CreateFactorV2(n *Node) *FactorV2 {
	factor := new(FactorV2)
	factor.CPT = n.CPT
	factor.Card = append(factor.Card, n.Numvalues)
	factor.Scope = append(factor.Scope, n)
	factor.Scope = append(factor.Scope, n.Parents...)
	factor.Strides = map[*Node]int{}
	factor.Strides[n] = 1
	s := n.Numvalues
	for _, v := range n.Parents {
		factor.Strides[v] = s
		s *= v.Numvalues
	}
	return factor
}

// FactorProduct permits to calculate Ψ(R, T, L), consider the net (R)->(T)->(L)
// We have to calculate multiple joins, We have P(R), P(T|R) and P(L|T), corrispondig to factors Φ(R) and Φ(R, T)
// First we calculate P(R,T) = Φ(R)Φ(R, T), so we have a factor net (R, T)->(L)
// Then we have P(L|R, T), to calculate the factor we do the same way and we have Ψ(R, T, L)
// This is for a vectorized version. deprecated.
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
	// fmt.Println("MAP", mapA, mapB, C)
	values := make([]int, 0)
	for i := 0; i < len(C.Values); i++ {
		values = append(values, i)
	}
	assignments := IndexToAssignment(values, C.Numvalues)
	// fmt.Println("assignmets", assignments)
	AssignmentToIndex(assignments, A.Numvalues)

}

// IndexToAssignment calculate var indexes for multiply factor. deprecated.
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
}

// AssignmentToIndex get var indexes of assignments in factor. deprecated.
func AssignmentToIndex(assignments *mu.Matrix, factorcard []int) {
	cardvect := mu.NewVector(factorcard)
	fmt.Println("CARDVECT", cardvect)
	//cprodD := mu.Cumprod(1, cardvect)
	if assignments.M == 1 || assignments.N == 1 {
		tmp := mu.Cumprod(mu.CreateNewSlice([]int{1}, factorcard[:len(factorcard)-1]))
		fmt.Println(tmp)
		// I = cumprod([1, D(1:end - 1)]) * (A(:) - 1) + 1
	} else {
		// I = sum(repmat(cumprod([1, D(1:end - 1)]), size(A, 1), 1) .* (A - 1), 2) + 1
	}

}

func intersectV2(A *FactorV2, B *FactorV2) []*Node {
	intersection := make([]*Node, 0)
	for _, node := range A.Scope {
		if containsNodeV2(B, node) {
			intersection = append(intersection, node)
		}
	}
	return intersection
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

func containsNodeV2(f *FactorV2, n *Node) bool {
	for _, node := range f.Scope {
		if node == n {
			return true
		}
	}
	return false
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
		// fmt.Println(variable)
		for index, cvar := range C.Var {
			if cvar == variable {
				pos = append(pos, index)
			}
		}
	}
	return pos
}

// calculate stride of variabile in var list, see p.359 Probabilist Graphical Models. Koller, Friedman.
func stride(v *Node, vars []*Node, cardinalities []int) int {
	var stride int
	var found bool = false
	for _, variable := range vars {
		if v == variable {
			found = true
		}
	}
	if !found {
		stride = 0
	} else {
		stride = 1
	}
	// for i, j := 0, len(cardinalities)-1; i < j; i, j = i+1, j-1 {
	// 	cardinalities[i], cardinalities[j] = cardinalities[j], cardinalities[i]
	// }
	for i, cardinality := range cardinalities {
		if v == vars[i] {
			break
		}
		stride *= cardinality
	}
	return stride
}

// MultiplyFactor takes in input two factor and multiply them, see p.359 Probabilist Graphical Models. Koller, Friedman. Algorithm 10.A.1
func MultiplyFactor(phi1 *FactorV2, phi2 *FactorV2, retcopy bool) *FactorV2 {
	variables := phi1.Scope
	for _, v := range phi2.Scope {
		if containvar(v, variables) == false {
			variables = append(variables, v)
		}
	}
	cardinality := phi1.Card
	for _, v := range phi2.Scope {
		if containvar(v, phi1.Scope) == false {
			cardinality = append(cardinality, v.Numvalues)
		}
	}
	var quantityvalues int = 1
	for i := 0; i < len(cardinality); i++ {
		quantityvalues *= cardinality[i]
	}
	values := make([]float64, quantityvalues)
	// Algorithm start
	j := 0
	k := 0
	assignment := make(map[*Node]int)
	for _, v := range variables {
		assignment[v] = 0
	}
	for i := 0; i < quantityvalues; i++ { // for l = 0 ... |X1 U X2|
		index := 0
		// fmt.Println("j", j, "k", k)
		for _, v := range variables { // for i = 0 ... |Var(X1 U X2)|
			index = index + assignment[v]*stride(v, variables, cardinality)
		}
		// fmt.Println("INDEX", index)
		values[index] = phi1.CPT[j] * phi2.CPT[k]
		// fmt.Println(values)
		for idx, variable := range variables {
			assignment[variable] = assignment[variable] + 1
			if assignment[variable] == cardinality[idx] {
				assignment[variable] = 0
				j = j - (cardinality[idx]-1)*phi1.stride(variable)
				k = k - (cardinality[idx]-1)*phi2.stride(variable)
			} else {
				j = j + phi1.stride(variable)
				k = k + phi2.stride(variable)
				break
			}
		}
	}
	psi := new(FactorV2)
	psi.Scope = variables
	psi.CPT = values
	psi.Card = cardinality
	psi.Strides = map[*Node]int{}
	for _, v := range psi.Scope {
		psi.Strides[v] = stride(v, psi.Scope, psi.Card)
	}
	if !retcopy {
		copier.Copy(phi1, psi)
		return nil
	}
	return psi
}

// DivideFactor takes in input two factor and divide them, see p.359 Probabilist Graphical Models. Koller, Friedman. Algorithm 10.A.1
// Value adjusted by following 9.3 p.297 Koller
func DivideFactor(phi1 *FactorV2, phi2 *FactorV2, retcopy bool) *FactorV2 {
	variables := phi1.Scope
	for _, v := range phi2.Scope {
		if containvar(v, variables) == false {
			variables = append(variables, v)
		}
	}
	cardinality := phi1.Card
	for _, v := range phi2.Scope {
		if containvar(v, phi1.Scope) == false {
			cardinality = append(cardinality, v.Numvalues)
		}
	}
	var quantityvalues int = 1
	for i := 0; i < len(cardinality); i++ {
		quantityvalues *= cardinality[i]
	}
	values := make([]float64, quantityvalues)
	// Algorithm start
	j := 0
	k := 0
	assignment := make(map[*Node]int)
	for _, v := range variables {
		assignment[v] = 0
	}
	for i := 0; i < quantityvalues; i++ { // for l = 0 ... |X1 U X2|
		index := 0
		// fmt.Println("j", j, "k", k)
		for _, v := range variables { // for i = 0 ... |Var(X1 U X2)|
			index = index + assignment[v]*stride(v, variables, cardinality)
		}
		if phi2.CPT[k] == 0 {
			values[index] = 0
		} else {
			values[index] = phi1.CPT[j] / phi2.CPT[k]
		}
		// fmt.Println(phi1.CPT[j], "/", phi2.CPT[k], "=", phi1.CPT[j]/phi2.CPT[k])
		// fmt.Println(values)
		for idx, variable := range variables {
			assignment[variable] = assignment[variable] + 1
			if assignment[variable] == cardinality[idx] {
				assignment[variable] = 0
				j = j - (cardinality[idx]-1)*phi1.stride(variable)
				k = k - (cardinality[idx]-1)*phi2.stride(variable)
			} else {
				j = j + phi1.stride(variable)
				k = k + phi2.stride(variable)
				break
			}
		}
	}
	psi := new(FactorV2)
	psi.Scope = variables
	psi.CPT = values
	psi.Card = cardinality
	psi.Strides = map[*Node]int{}
	for _, v := range psi.Scope {
		psi.Strides[v] = stride(v, psi.Scope, psi.Card)
	}
	if !retcopy {
		copier.Copy(phi1, psi)
		return nil
	}
	fmt.Println("phi1", phi1)
	fmt.Println("phi2", phi2)
	tmp := make([]float64, quantityvalues)

	valdivide := quantityvalues / len(phi2.CPT)
	s := 0
	for i := 0; i < len(phi2.CPT); i++ {
		for j := s; j < s+valdivide; j++ {
			// fmt.Println(phi1.CPT[i+j], " / ", phi2.CPT[i], "=", phi1.CPT[i+j]/phi2.CPT[i])
			if phi2.CPT[i] == 0 {
				tmp[i+j] = 0
			} else {
				tmp[i+j] = phi1.CPT[i+j] / phi2.CPT[i]
			}
		}
		s += valdivide - 1
	}
	fmt.Println("tmp", tmp)
	psi.CPT = tmp
	return psi
}

func containvar(target *Node, list []*Node) bool {
	found := false
	for _, variable := range list {
		if variable == target {
			found = true
		}
	}
	return found
}

func (f *FactorV2) stride(variable *Node) int {
	found := false
	for v := range f.Strides {
		if v == variable {
			found = true
		}
	}
	if f.Strides == nil {
		f.Strides = make(map[*Node]int)
	}
	if !found {
		f.Strides[variable] = stride(variable, f.Scope, f.Card)
	}
	return f.Strides[variable]
}

// Marginalize one ore more variables from a factor, i want to return a factor with summed out variables in input
func (f *FactorV2) Marginalize(retcopy bool, variables ...*Node) *FactorV2 {
	fmt.Println("Marginalizing ")
	for i, variable := range variables {
		fmt.Println(i, ".", variable.Name)
		if !containvar(variable, f.Scope) {
			fmt.Errorf("Variable not in scope")
		}
	}
	phi := new(FactorV2)

	varindex := make([]int, 0)
	indextokeep := make([]int, 0)

	// get indexes of variables to remove
	for ind, variable := range f.Scope {
		// fmt.Println("STRIDE", f.stride(variable))
		if containvar(variable, variables) {
			varindex = append(varindex, ind)
		} else {
			indextokeep = append(indextokeep, ind)
		}
	}

	// for every variable not to remove, not contained in
	for index, variable := range f.Scope {
		if in, _ := contains(variables, variable); in == false {
			phi.Scope = append(phi.Scope, variable)
			phi.Card = append(phi.Card, f.Card[index])
		}
	}
	// fmt.Println("Var indexs", varindex)
	fmt.Println("Index to keep", indextokeep)
	fmt.Println("card", phi.Card)
	// fmt.Println("NEW", phi)
	phi.Strides = map[*Node]int{}
	for _, v := range phi.Scope {
		phi.Strides[v] = 1 // f.stride(v)
	}
	fmt.Println("PHI", phi)
	nvalues := 1
	for _, card := range phi.Card {
		nvalues *= card
	}

	//phi.CPT = make([]float64, nvalues)
	tmp := make([]float64, nvalues)
	strides := make(map[*Node]int)
	newindex := make([]int, 0)
	for _, v := range variables {
		strides[v] = f.Strides[v]
		newindex = append(newindex, f.Strides[v]-1)
	}
	color.HiRed("\nstrides\n", strides)
	color.HiRed("\nnewindex\n", newindex)
	for v, s := range strides {
		fmt.Println(v, s)
		for i := 0; i < nvalues; i++ {
			for j := i; j < s*nvalues; j += s {
				tmp[i] += f.CPT[j]
				fmt.Println(tmp, i)
			}
		}
	}
	fmt.Println(f.CPT)
	fmt.Println("MY TMP", tmp)
	tensor := ts.New(ts.WithBacking(f.CPT), ts.WithShape(f.Card...))
	fmt.Println(tensor)
	fmt.Println(varindex)
	// for i, v := range varindex {
	// 	if v == 0 {
	// 		varindex[i] = 1
	// 	}
	// 	if v == 1 {
	// 		varindex[i] = 0
	// 	}
	// }
	summed, _ := ts.Sum(tensor, newindex...) // 0 sum on x, 1 sum on y, 2 sum on z ecc..
	fmt.Println(summed)
	// phi.CPT = summed.Data().([]float64)
	phi.CPT = tmp
	if !retcopy {
		copier.Copy(f, phi)
		return nil
	}
	fmt.Println("NEW phi", phi)
	return phi

}

// Normalize factor CPT
func (f *FactorV2) Normalize() {
	var sum float64 = 0
	for _, value := range f.CPT {
		sum += value
	}
	for i, value := range f.CPT {
		f.CPT[i] = value / sum
	}
}
