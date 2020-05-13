package main

import (
	"fmt"
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
	bn.FactorProduct(ifact, hfact)
}
