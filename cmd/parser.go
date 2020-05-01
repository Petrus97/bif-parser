package cmd

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// Probabilities assigned in the model
type Probabilities struct {
	states map[string]float64
}

// Node if net
type Node struct {
	Name      string
	Type      string
	numvalues int
	prob      Probabilities
	CPT       [][]float64
	parents   []*Node
	child     []*Node
}

// BN bayiesian Network
type BN struct {
	Nodes []*Node
}

func (n *Node) addParents(nodes ...*Node) {
	for i := 0; i < len(nodes); i++ {
		n.parents = append(n.parents, nodes[i])
		nodes[i].child = append(nodes[i].child, n)
	}
}

// not efficient but simple, for bigger net is better use an hash list
func (bn *BN) getNode(name string) *Node {
	for i := 0; i < len(bn.Nodes); i++ {
		if bn.Nodes[i].Name == name {
			return bn.Nodes[i]
		}
	}
	return nil
}

func (bn *BN) updatePrior(matchprior map[string]string) {
	node := bn.getNode(matchprior["var"])
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

func (bn *BN) createNode(match map[string]string) {
	probs := strings.Split(match["state"], ", ")
	n, _ := strconv.Atoi(match["nval"])
	states := map[string]float64{}
	for i := 0; i < n; i++ {
		states[probs[i]] = 0
	}
	node := Node{
		Name:      match["var"],
		Type:      match["type"],
		numvalues: n,
		prob: Probabilities{
			states: states,
		},
		CPT:     make([][]float64, 0),
		parents: make([]*Node, 0),
		child:   make([]*Node, 0),
	}
	bn.Nodes = append(bn.Nodes, &node)
}

func (bn *BN) listNodes() {
	fmt.Println("Listing nodes")
	for _, node := range bn.Nodes {
		fmt.Println("Name: ", node.Name)
		fmt.Println("type: ", node.Type)
		fmt.Println("nval: ", node.numvalues)
		fmt.Println("CPT: ", node.CPT)
		fmt.Println("states: ", node.prob.states)
		fmt.Println("child: ", node.child)
		fmt.Println("parents: ", node.parents)
		fmt.Println("")
	}
}

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadBIF function that takes in input path to
// .bif file ancd create bayesian network
func ReadBIF(filepath string) *BN {
	bn := BN{
		Nodes: make([]*Node, 0),
	}

	varpattern, _ := regexp.Compile("variable (?P<var>[[:graph:]]*) \\{\n  type (?P<type>[a-z]*) \\[ (?P<nval>\\d+) \\] \\{ (?P<state>.*) \\};")
	priorprobpattern, _ := regexp.Compile("probability \\( (?P<var>[^|,]+) \\) \\{\n  table (?P<prior>.+);")
	condprobpattern, _ := regexp.Compile("probability \\( (?P<child>.+) \\| (?P<parents>.+) \\) \\{\n")
	condprobpattern2, _ := regexp.Compile("  \\((?P<key>.+)\\) (?P<values>.+);")

	file, err := ioutil.ReadFile(filepath)
	errorCheck(err)

	matchvar := map[string]string{}
	matchprior := map[string]string{}
	matchparents := map[string]string{}

	variables := varpattern.FindAllStringSubmatch(string(file), -1)
	cpts := condprobpattern2.FindAllStringSubmatch(string(file), -1)

	keys := varpattern.SubexpNames()
	for _, v := range variables { // for every variable
		for i, mName := range v { // i is the index, mName is the matched name
			matchvar[keys[i]] = mName
		}
		bn.createNode(matchvar)
	}
	fmt.Println(len(bn.Nodes))
	keys = priorprobpattern.SubexpNames()
	for _, p := range priorprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach prior probability found
		for i, mName := range p { // i is the index, mName is the matched name
			matchprior[keys[i]] = mName
		}
		bn.updatePrior(matchprior)
	}
	j := 0
	keys = condprobpattern.SubexpNames()
	for _, rel := range condprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach relation child parents
		for i, mName := range rel { // i is the index, mName is the matched name
			matchparents[keys[i]] = mName
		}
		node := bn.getNode(matchparents["child"])
		parentnames := strings.Split(matchparents["parents"], ", ")
		nval := 1
		for i := 0; i < len(parentnames); i++ {
			node.addParents(bn.getNode(parentnames[i]))
			nval *= node.parents[i].numvalues
		}
		// fmt.Println(node.Name, nval)
		for i := 0; i < nval; i++ {
			if i >= nval {
				break
			}
			values := make([]float64, 0)
			svalues := strings.Split(cpts[i+j][2], ", ")
			for _, s := range svalues {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					values = append(values, f)
				}
			}
			node.CPT = append(node.CPT, values)
		}
		j += nval
	}
	// bn.listNodes()
	return &bn
}
