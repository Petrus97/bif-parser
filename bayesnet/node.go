package bayesnet

import "fmt"

// Probabilities assigned in the model
type Probabilities struct {
	States []string
	Prob   []float64
}

// Node if net
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

func (n *Node) GetPotential() {
	pot := make([]float64, 0)
	k := 0
	for _, par := range n.Parents {
		for j := 0; j < par.Numvalues; j++ {
			for i := k; i < k+n.Numvalues; i++ {
				fmt.Println(par.CPT[j], n.CPT[i], "=", par.CPT[j]*n.CPT[i])
				pot = append(pot, par.CPT[j]*n.CPT[i])
			}
			k += par.Numvalues
		}
	}
	fmt.Println(pot)
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
