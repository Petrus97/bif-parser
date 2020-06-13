package bayesnet

// Probabilities assigned in the model
type Probabilities struct {
	States []string
	Prob   []float64
}

// Node of the net
type Node struct {
	Name      string
	Type      string
	Numvalues int
	Prob      Probabilities
	CPT       []float64
	Domain    [][]string
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

func difference(alist []*Node, blist []*Node) []*Node {
	// A - B
	clist := make([]*Node, 0)
	for _, n := range alist {
		if ok := containvar(n, blist); ok == false {
			clist = append(clist, n)
		}
	}
	for _, n := range blist {
		if ok := containvar(n, alist); ok == false {
			clist = append(clist, n)
		}
	}
	return clist

}
