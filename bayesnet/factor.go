package bayesnet

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

}

func intersect(A *Factor, B *Factor) {

}
