package parser

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"

	bn "github.com/Petrus97/bif-parser/bayesnet"
)

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

// ReadBIF function that takes in input path to
// .bif file ancd create bayesian network
func ReadBIF(filepath string) *bn.BN {
	bn := bn.BN{
		Nodes: make([]*bn.Node, 0),
	}

	varpattern, _ := regexp.Compile("variable (?P<var>[[:graph:]]*) \\{\n  type (?P<type>[a-z]*) \\[ (?P<nval>\\d+) \\] \\{ (?P<state>.*) \\};")
	priorprobpattern, _ := regexp.Compile("probability \\( (?P<var>[^|,]+) \\) \\{\n  table (?P<prior>.+);")
	condprobpattern, _ := regexp.Compile("probability \\( (?P<child>.+) \\| (?P<parents>.+) \\) \\{\n")
	condprobpattern2, _ := regexp.Compile("  \\((?P<domain>.+)\\) (?P<values>.+);")

	file, err := ioutil.ReadFile(filepath)
	errorCheck(err)

	matchvar := map[string]string{}
	matchprior := map[string]string{}
	matchparents := map[string]string{}

	variables := varpattern.FindAllStringSubmatch(string(file), -1)
	cpts := condprobpattern2.FindAllStringSubmatch(string(file), -1)
	// otherkey := condprobpattern2.SubexpNames()
	// fmt.Println(otherkey)
	//fmt.Println(condprobpattern2.FindAllStringSubmatch(string(file), -1))
	keys := varpattern.SubexpNames()
	for _, v := range variables { // for every variable
		for i, mName := range v { // i is the index, mName is the matched name
			matchvar[keys[i]] = mName
		}
		bn.CreateNode(matchvar)
	}
	//fmt.Println(len(bn.Nodes))
	keys = priorprobpattern.SubexpNames()
	for _, p := range priorprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach prior probability found
		for i, mName := range p { // i is the index, mName is the matched name
			matchprior[keys[i]] = mName
		}
		bn.UpdatePrior(matchprior)
	}
	j := 0
	keys = condprobpattern.SubexpNames()
	for _, rel := range condprobpattern.FindAllStringSubmatch(string(file), -1) { // foreach relation child parents
		for i, mName := range rel { // i is the index, mName is the matched name
			matchparents[keys[i]] = mName
		}
		node := bn.GetNode(matchparents["child"])
		parentnames := strings.Split(matchparents["parents"], ", ")
		nval := 1
		for i := 0; i < len(parentnames); i++ {
			node.AddParents(bn.GetNode(parentnames[i]))
			nval *= node.Parents[i].Numvalues
		}
		node.Domain = make([][]string, 0)
		// fmt.Println(node.Name, nval)
		for i := 0; i < nval; i++ {
			if i >= nval {
				break
			}
			// values := make([]float64, 0)
			svalues := strings.Split(cpts[i+j][2], ", ")
			sdomains := strings.Split(cpts[i+j][1], ", ")
			// fmt.Println(sdomains)
			for _, s := range svalues {
				if f, err := strconv.ParseFloat(s, 64); err == nil {
					node.CPT = append(node.CPT, f)
				}
			}
			for t := 0; t < node.Numvalues; t++ {
				myDomain := make([]string, 0)
				myDomain = append(myDomain, node.Prob.States[t])
				for _, s := range sdomains {
					myDomain = append(myDomain, s)
				}
				node.Domain = append(node.Domain, myDomain)
			}
		}
		j += nval
	}
	return &bn
}
