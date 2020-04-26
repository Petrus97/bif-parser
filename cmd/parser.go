package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
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
	CPT       map[int]float64
	parents   []*Node
	child     []*Node
}

// BN bayiesian Network
type BN struct {
	nodes []*Node
}

func (n *Node) addParents(nodes ...*Node) {
	for i := 0; i < len(nodes); i++ {
		n.parents = append(n.parents, nodes[i])
		nodes[i].child = append(nodes[i].child, n)
	}
}

// not efficient but simple, for bigger net is better use an hash list
func (bn *BN) getNode(name string) *Node {
	for i := 0; i < len(bn.nodes); i++ {
		if bn.nodes[i].Name == name {
			return bn.nodes[i]
		}
	}
	return nil
}

func (bn *BN) updatePrior(matchprior map[string]string) {
	node := bn.getNode(matchprior["var"])
	probabilities := strings.Split(matchprior["prior"], ", ")
	i := 0
	for name := range node.prob.states {
		node.prob.states[name], _ = strconv.ParseFloat(probabilities[i], 64)
		node.CPT[i] = node.prob.states[name]
		i++
	}
}

func (n *Node) updateCPT(cptvalues string, index int) {
	splitted := strings.Split(cptvalues, ", ")
	fmt.Println(splitted)
	for i := 0; i < len(splitted); i++ {
		value, _ := strconv.ParseFloat(splitted[i], 64)
		n.CPT[index] = value
		index++
	}

}

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func ReadBIF_V2() {
	bn := BN{
		nodes: make([]*Node, 0),
	}
	varpattern, _ := regexp.Compile("variable (?P<var>[[:alpha:]]*) \\{\n  type (?P<type>[a-z]*) \\[ (?P<nval>\\d+) \\] \\{ (?P<state>.*) \\};")
	priorprobpattern, _ := regexp.Compile("probability \\( (?P<var>[^|,]+) \\) \\{\n  table (?P<prior>.+);")
	condprobpattern, _ := regexp.Compile("probability \\( (?P<child>.+) \\| (?P<parents>.+) \\) \\{\n")
	condprobpattern2, _ := regexp.Compile("  \\((?P<key>.+)\\) (?P<values>.+);")
	file, err := ioutil.ReadFile("data/survey.bif")
	errorCheck(err)
	matchvar := map[string]string{}
	matchprior := map[string]string{}
	matchparents := map[string]string{}
	// variables := make([]string, 0)
	// variables = append(variables, varpattern.FindAllString(string(file), -1)...)
	// fmt.Println(varpattern.FindAllString(string(file), -1))
	variables := varpattern.FindAllStringSubmatch(string(file), -1)
	cpts := condprobpattern2.FindAllStringSubmatch(string(file), -1)
	fmt.Println("COND:\n", condprobpattern.FindAllStringSubmatch(string(file), -1))
	fmt.Println(condprobpattern.SubexpNames())

	fmt.Println("COND2:\n", condprobpattern2.FindAllStringSubmatch(string(file), -1))
	fmt.Println(condprobpattern2.SubexpNames())
	// variables := varpattern.FindAllString(string(file), -1)
	keys := varpattern.SubexpNames()
	for _, v := range variables { // for every variable
		for i, mName := range v { // i is the index, mName is the matched name
			matchvar[keys[i]] = mName
		}
		bn.createNode(matchvar)
	}
	// TODO make some clean in the code
	keys = priorprobpattern.SubexpNames()
	for _, p := range priorprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach prior probability found
		// fmt.Println(p)
		for i, mName := range p { // i is the index, mName is the matched name
			matchprior[keys[i]] = mName
		}
		bn.updatePrior(matchprior)
	}
	j := 0
	keys = condprobpattern.SubexpNames()
	for _, rel := range condprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach relation child parents
		// fmt.Println(rel)
		for i, mName := range rel { // i is the index, mName is the matched name
			matchparents[keys[i]] = mName
		}
		node := bn.getNode(matchparents["child"])
		parentnames := strings.Split(matchparents["parents"], ", ")
		for _, p := range parentnames {
			node.addParents(bn.getNode(p))
		}
		tabledim := 0
		for _, p := range node.parents {
			tabledim += p.numvalues
		}
		for i := 0; i < node.numvalues*tabledim; i++ {
			node.updateCPT(cpts[i+j][2], len(cpts[i+j][2]))
		}
		j += tabledim
	}
	bn.listNodes()
}

// ReadBIF is use to read a bif file
func ReadBIF() {
	// bn := BN{
	// 	nodes: make([]*Node, 0),
	// }
	variable, _ := regexp.Compile("^variable\\s[[:alpha:]].*\\{$")
	varpattern, _ := regexp.Compile("[^variable\\s][[:alpha:]]*")
	typepattern, _ := regexp.Compile("  type discrete \\[ \\d+ \\] \\{ (.+) \\};\\s*")
	priorprobpattern, _ := regexp.Compile("probability \\( ([^|]+) \\) \\{\\s*")
	priorprobpattern2, _ := regexp.Compile("  table (.+);\\s*")
	condprobpattern, _ := regexp.Compile("probability \\( (.+) \\| (.+) \\) \\{\\s*")
	condprobpattern2, _ := regexp.Compile("  \\((.+)\\) (.+);\\s*")
	//remvar, _ := regexp.Compile("variable*{")
	nval, _ := regexp.Compile("[^  type discrete \\[ ]+")

	file, err := os.Open("data/earthquake.bif")
	errorCheck(err)
	defer file.Close()
	variables := make([]string, 0)
	types := make([]string, 0)

	// reader := bufio.NewReader(os.Stdin)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if s := variable.FindString(scanner.Text()); s != "" {
			variables = append(variables, varpattern.FindString(scanner.Text()))
		}
		if t := typepattern.FindString(scanner.Text()); t != "" {
			fmt.Println(t)
			fmt.Println(nval.FindString(t))
			types = append(types, t)
		}
		if p1 := priorprobpattern.FindString(scanner.Text()); p1 != "" {
			fmt.Println(p1)
		}
		if p2 := priorprobpattern2.FindString(scanner.Text()); p2 != "" {
			fmt.Println(p2)
		}
		if c1 := condprobpattern.FindString(scanner.Text()); c1 != "" {
			fmt.Println(c1)
		}
		if c2 := condprobpattern2.FindString(scanner.Text()); c2 != "" {
			fmt.Println(c2)
		}

	}
	// r, _ := regexp.Compile("[^variable\\s][[:alpha:]]*")
	// for i := 0; i < len(variables); i++ {
	// 	bn.createNode(variables[i])
	// }
	// bn.listNodes()
}

func (bn *BN) createNode(match map[string]string) {
	probs := strings.Split(match["state"], ", ")
	// fmt.Println(probs)
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
		CPT:     make(map[int]float64),
		parents: make([]*Node, 0),
		child:   make([]*Node, 0),
	}
	bn.nodes = append(bn.nodes, &node)
}

func (bn *BN) listNodes() {
	fmt.Println("Listing nodes")
	for _, node := range bn.nodes {
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
