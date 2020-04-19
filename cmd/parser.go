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

func (b *BN) getNode(name string) *Node {
	for i := 0; i < len(b.nodes); i++ {
		if b.nodes[i].Name == name {
			return b.nodes[i]
		}
	}
	return nil
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
	priorprobpattern, _ := regexp.Compile("probability \\( (?P<var>[^|].+) \\) \\{\n  table (?P<prior>.*);")
	condprobpattern, _ := regexp.Compile("probability \\( (?P<child>.+) \\| (?P<parents>.+) \\) \\{\n")
	condprobpattern2, _ := regexp.Compile("  \\((?P<key>.+)\\) (?P<values>.+);")
	file, err := ioutil.ReadFile("data/earthquake.bif")
	// fmt.Println(string(file))
	errorCheck(err)
	matchvar := map[string]string{}
	matchprior := map[string]string{}
	// variables := make([]string, 0)
	// variables = append(variables, varpattern.FindAllString(string(file), -1)...)
	// fmt.Println(varpattern.FindAllString(string(file), -1))
	variables := varpattern.FindAllStringSubmatch(string(file), -1)
	fmt.Println("PRIOR:\n", priorprobpattern.FindAllStringSubmatch(string(file), -1))
	fmt.Println(priorprobpattern.SubexpNames())

	fmt.Println("COND:\n", condprobpattern.FindAllStringSubmatch(string(file), -1))
	fmt.Println(condprobpattern.SubexpNames())

	fmt.Println("COND2:\n", condprobpattern2.FindAllStringSubmatch(string(file), -1))
	fmt.Println(condprobpattern2.SubexpNames())
	// variables := varpattern.FindAllString(string(file), -1)
	keys := varpattern.SubexpNames()
	// fmt.Printf("match: '%s'\nname: '%s'\n", variables[0], keys[1])
	fmt.Println(keys)
	for _, v := range variables { // for every variable
		for i, n := range v { // i is the index, n is the matched name
			// fmt.Printf("%d. match='%s'\tname='%s'\n", i, n, keys[i])
			matchvar[keys[i]] = n
		}
		bn.createNode(matchvar)
	}
	keys = priorprobpattern.SubexpNames()
	// TODO make some clean in the code
	for _, p := range priorprobpattern.FindAllStringSubmatch(string(file), -1) {
		fmt.Println(p)
		for i, n := range p {
			matchprior[keys[i]] = n
		}
		// TODO updatePrior(map[string]string)
		// TODO updateCPT()
	}
	for key, val := range matchprior {
		fmt.Printf("key: '%s'\t-\tvalue: '%s'\n", key, val)
	}
	fmt.Println("LEN:", len(variables))
	// for key, val := range matchvar {
	// 	fmt.Printf("key: '%s'\t-\tvalue: '%s'\n", key, val)
	// }
	bn.listNodes()
	// for _, v := range variables {
	// 	fmt.Println(v)
	// }
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
	}
	bn.nodes = append(bn.nodes, &node)
}

func (bn *BN) listNodes() {
	fmt.Println("Listing nodes")
	for _, node := range bn.nodes {
		fmt.Println(node.Name)
		fmt.Println(node.Type)
		fmt.Println(node.numvalues)
		fmt.Println(node.prob.states)
		fmt.Println("")
	}
}
