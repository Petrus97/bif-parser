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
	for name := range node.Prob.States {
		if f, err := strconv.ParseFloat(probabilities[i], 64); err == nil {
			node.Prob.Prob = append(node.Prob.Prob, f)
			node.CPT = append(node.CPT, node.Prob.Prob[name])
			i++
		}
	}
	node.Domain = make([][]string, len(node.Prob.States))
	for i := 0; i < len(node.Prob.States); i++ {
		node.Domain[i] = append(node.Domain[i], node.Prob.States[i])
	}

}

// CreateNode takes in input a map to create a node
func (bn *BN) CreateNode(match map[string]string) {
	probs := strings.Split(match["state"], ", ")
	n, _ := strconv.Atoi(match["nval"])
	states := make([]string, 0)
	domains := map[string][]string{}
	for i := 0; i < n; i++ {
		states = append(states, probs[i])
		domains[probs[i]] = nil
	}
	node := Node{
		Name:      match["var"],
		Type:      match["type"],
		Numvalues: n,
		Prob: Probabilities{
			States: states,
			Prob:   make([]float64, 0),
		},
		CPT:     make([]float64, 0),
		Domain:  make([][]string, 0),
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
		fmt.Println("Domain: ", node.Domain)
		fmt.Println("States: ", node.Prob.States)
		fmt.Println("Prob: ", node.Prob.Prob)
		fmt.Println("child: ", node.Child)
		fmt.Println("parents: ", node.Parents)
		fmt.Println("")
	}
}
