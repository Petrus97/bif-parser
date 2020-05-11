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
	watson := myNet.GetNode("Watson")
	holmes := myNet.GetNode("Holmes")
	// fmt.Println("Watson")
	// watson.GetPotential()
	// holmes := myNet.GetNode("Holmes")
	// fmt.Println("HOLMES")
	// holmes.GetPotential()
	// myNet.MatchDomains()
	wfact := bn.CreateFactor(watson)
	hfact := bn.CreateFactor(holmes)
	fmt.Println(wfact)
	fmt.Println(hfact)
}
