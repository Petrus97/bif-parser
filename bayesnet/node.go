package bayesnet

// Probabilities assigned in the model
type Probabilities struct {
	states map[string]float64
}

// Node if net
type Node struct {
	Name      string
	Type      string
	Numvalues int
	prob      Probabilities
	CPT       [][]float64
	Parents   []*Node
	Child     []*Node
}

// AddParents takes a lisit of nodes to add in Parents list, update also child
func (n *Node) AddParents(nodes ...*Node) {
	for i := 0; i < len(nodes); i++ {
		n.Parents = append(n.Parents, nodes[i])
		nodes[i].Child = append(nodes[i].Child, n)
	}
}
