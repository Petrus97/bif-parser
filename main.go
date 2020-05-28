package main

import (
	"fmt"
	"strings"
	"time"

	bn "github.com/Petrus97/bif-parser/bayesnet"
	"github.com/Petrus97/bif-parser/parser"
)

func main() {
	fmt.Println("Hello")
	start := time.Now()
	myNet := parser.ReadBIF("data/icy_roads.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
	icy := myNet.GetNode("Icy")
	holmes := myNet.GetNode("Holmes")

	// Naive method to create join prob table icy-holmes
	var jpt []float64
	if len(holmes.CPT) > len(icy.CPT) {
		jpt = make([]float64, len(holmes.CPT))
		for i := 0; i < len(holmes.Domain); i++ {
			jpt[i] = holmes.CPT[i]
			fmt.Println(holmes.CPT[i], "\t", holmes.Domain[i])
			for j := 1; j < len(holmes.Domain[i]); j++ {
				dom := holmes.Domain[i][j]
				parent := holmes.Parents[j-1]
				var index int
				for k := 0; k < len(parent.Prob.States); k++ {
					if strings.Compare(dom, parent.Prob.States[k]) == 0 {
						index = k
					}
				}
				jpt[i] *= parent.CPT[index]
			}
		}
	}
	fmt.Println(jpt)
	// fmt.Println("Watson")
	// watson.GetPotential()
	// holmes := myNet.GetNode("Holmes")
	// fmt.Println("HOLMES")
	// holmes.GetPotential()
	// myNet.MatchDomains()
	ifact := bn.CreateFactor(icy)
	hfact := bn.CreateFactor(holmes)
	fmt.Println(ifact)
	fmt.Println(hfact)
	// bn.FactorProduct(ifact, hfact)

	icyfact := new(bn.FactorV2)
	icyfact.CPT = icy.CPT
	icyfact.Scope = append(icyfact.Scope, icy)
	icyfact.Card = append(icyfact.Card, icy.Numvalues)

	holmesfact := new(bn.FactorV2)
	holmesfact.CPT = holmes.CPT
	holmesfact.Scope = append(holmesfact.Scope, holmes)
	holmesfact.Scope = append(holmesfact.Scope, holmes.Parents...)
	holmesfact.Card = append(holmesfact.Card, holmes.Numvalues)
	for _, p := range holmes.Parents {
		holmesfact.Card = append(holmesfact.Card, p.Numvalues)
	}

	fmt.Println("F V2 Icy", icyfact)
	fmt.Println("F V2 Holmes", holmesfact)

	ih := bn.MultiplyFactor(icyfact, holmesfact)
	fmt.Println("F V2 Icy", icyfact)
	fmt.Println("F V2 Holmes", holmesfact)
	fmt.Println("Icy-Holmes", ih)

	ih.Marginalize(false, holmes)
	// h := bn.DivideFactor(&ih, holmesfact)
	// fmt.Println(h)
}
