package bayesnet

import (
	"fmt"
	"strconv"
	"strings"
)

// BN bayiesian Network
type BN struct {
	Nodes []*Node
}

// GetNode of the network by name; not efficient but simple, for bigger net is better use an hash list
func (bn *BN) GetNode(name string) *Node {
	for i := 0; i < len(bn.Nodes); i++ {
		if bn.Nodes[i].Name == name {
			return bn.Nodes[i]
		}
	}
	return nil
}

// UpdatePrior probabilities during parsing
func (bn *BN) UpdatePrior(matchprior map[string]string) {
	node := bn.GetNode(matchprior["var"])
	probabilities := strings.Split(matchprior["prior"], ", ")
	i := 0
	var values []float64
	values = make([]float64, 0)
	for name := range node.prob.states {
		node.prob.states[name], _ = strconv.ParseFloat(probabilities[i], 64)
		values = append(values, node.prob.states[name])
		i++
	}
	node.CPT = append(node.CPT, values)
}

// CreateNode takes in input a map to create a node
func (bn *BN) CreateNode(match map[string]string) {
	probs := strings.Split(match["state"], ", ")
	n, _ := strconv.Atoi(match["nval"])
	states := map[string]float64{}
	for i := 0; i < n; i++ {
		states[probs[i]] = 0
	}
	node := Node{
		Name:      match["var"],
		Type:      match["type"],
		Numvalues: n,
		prob: Probabilities{
			states: states,
		},
		CPT:     make([][]float64, 0),
		Parents: make([]*Node, 0),
		Child:   make([]*Node, 0),
	}
	bn.Nodes = append(bn.Nodes, &node)
}

// ListNodes of network
func (bn *BN) ListNodes() {
	fmt.Println("Listing nodes")
	for _, node := range bn.Nodes {
		fmt.Println("Name: ", node.Name)
		fmt.Println("type: ", node.Type)
		fmt.Println("nval: ", node.Numvalues)
		fmt.Println("CPT: ", node.CPT)
		fmt.Println("states: ", node.prob.states)
		fmt.Println("child: ", node.Child)
		fmt.Println("parents: ", node.Parents)
		fmt.Println("")
	}
}
