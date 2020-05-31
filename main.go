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
	myNet := parser.ReadBIF("data/wet_grass.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
	rain := myNet.GetNode("Rain")
	sprinkler := myNet.GetNode("Sprinkler")
	holmes := myNet.GetNode("Holmes")
	watson := myNet.GetNode("Watson")

	rainfact := bn.CreateFactorV2(rain)
	sprinklerfact := bn.CreateFactorV2(sprinkler)
	holmesfact := bn.CreateFactorV2(holmes)
	watsonfact := bn.CreateFactorV2(watson)

	fmt.Println(rainfact)
	fmt.Println(sprinklerfact)
	fmt.Println(holmesfact)
	fmt.Println(watsonfact)

	wr := bn.MultiplyFactor(rainfact, watsonfact)
	fmt.Println("WR", wr)
	fmt.Println("rain", rainfact)
	fmt.Println("holmes", holmesfact)
	rs := bn.MultiplyFactor(sprinklerfact, rainfact)
	fmt.Println("RS", rs)
	hrs := bn.MultiplyFactor(&rs, holmesfact)

	fmt.Println("HRS", hrs)
	sepr1 := bn.DivideFactor(&wr, rainfact)
	fmt.Println("SEP_R1", sepr1)
	sper2 := bn.DivideFactor(&hrs, rainfact)
	fmt.Println("SEPR2", sper2)
	// hrs.Marginalize(false, holmes, sprinkler)

	// fmt.Println("WR", wr)

	// fmt.Println("Watson")
	// watson.GetPotential()
	// holmes := myNet.GetNode("Holmes")
	// fmt.Println("HOLMES")
	// holmes.GetPotential()
	// myNet.MatchDomains()
	// ifact := bn.CreateFactor(icy)
	// hfact := bn.CreateFactor(holmes)
	// fmt.Println(ifact)
	// fmt.Println(hfact)
	// bn.FactorProduct(ifact, hfact)

	// icy := myNet.GetNode("Icy")

	// icyfact := new(bn.FactorV2)
	// icyfact.CPT = icy.CPT
	// icyfact.Scope = append(icyfact.Scope, icy)
	// icyfact.Card = append(icyfact.Card, icy.Numvalues)

	// holmes := myNet.GetNode("Holmes")

	// holmesfact := new(bn.FactorV2)
	// holmesfact.CPT = holmes.CPT
	// holmesfact.Scope = append(holmesfact.Scope, holmes)
	// holmesfact.Scope = append(holmesfact.Scope, holmes.Parents...)
	// holmesfact.Card = append(holmesfact.Card, holmes.Numvalues)
	// for _, p := range holmes.Parents {
	// 	holmesfact.Card = append(holmesfact.Card, p.Numvalues)
	// }

	// watson := myNet.GetNode("Watson")
	// watsonfact := new(bn.FactorV2)
	// watsonfact.CPT = holmes.CPT
	// watsonfact.Scope = append(watsonfact.Scope, watson)
	// watsonfact.Scope = append(watsonfact.Scope, watson.Parents...)
	// watsonfact.Card = append(watsonfact.Card, watson.Numvalues)
	// for _, p := range watson.Parents {
	// 	watsonfact.Card = append(watsonfact.Card, p.Numvalues)
	// }

	// fmt.Println("F V2 Icy", icyfact)
	// fmt.Println("F V2 Holmes", holmesfact)
	// fmt.Println("F V2 Watson", watsonfact)

	// ih := bn.MultiplyFactor(icyfact, holmesfact)
	// fmt.Println("F V2 Icy", icyfact)
	// fmt.Println("F V2 Holmes", holmesfact)
	// fmt.Println("Icy-Holmes", ih)

	// iw := bn.MultiplyFactor(icyfact, watsonfact)
	// fmt.Println("F V2 Icy", icyfact)
	// fmt.Println("F V2 Watson", watsonfact)
	// fmt.Println("Icy-Watson", iw)

	// ih.Marginalize(false, icy)
	// iw.Marginalize(false, watson)
	// h := bn.DivideFactor(&ih, holmesfact)
	// fmt.Println(h)
}
