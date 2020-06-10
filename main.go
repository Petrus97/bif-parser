package main

import (
	"fmt"
	"time"

	bn "github.com/Petrus97/bif-parser/bayesnet"
	"github.com/Petrus97/bif-parser/parser"
)

func main() {
	start := time.Now()
	myNet := parser.ReadBIF("data/icy_roads.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
	// WET GRASS
	// rain := myNet.GetNode("Rain")
	// sprinkler := myNet.GetNode("Sprinkler")
	// holmes := myNet.GetNode("Holmes")
	// watson := myNet.GetNode("Watson")

	// rainfact := bn.CreateFactorV2(rain)
	// sprinklerfact := bn.CreateFactorV2(sprinkler)
	// holmesfact := bn.CreateFactorV2(holmes)
	// watsonfact := bn.CreateFactorV2(watson)

	// fmt.Println(rainfact)
	// fmt.Println(sprinklerfact)
	// fmt.Println(holmesfact)
	// fmt.Println(watsonfact)

	// wr := bn.MultiplyFactor(rainfact, watsonfact)
	// fmt.Println("WR", wr)
	// fmt.Println("rain", rainfact)
	// fmt.Println("holmes", holmesfact)
	// rs := bn.MultiplyFactor(sprinklerfact, rainfact)
	// fmt.Println("RS", rs)
	// hrs := bn.MultiplyFactor(&rs, holmesfact)
	// fmt.Println("HRS", hrs)
	// // sepr1 := bn.DivideFactor(&wr, rainfact)
	// // fmt.Println("SEP_R1", sepr1)
	// // sper2 := bn.DivideFactor(&hrs, rainfact)
	// // fmt.Println("SEPR2", sper2)
	// // hrs.Marginalize(false, holmes, sprinkler)
	// // holmesEv := bn.EnterEvidence(holmes, true)
	// // res := bn.MultiplyFactor()
	// // fmt.Println("WR", wr)

	// ICY ROADS
	icy := myNet.GetNode("Icy")

	icyfact := new(bn.FactorV2)
	icyfact.CPT = icy.CPT
	icyfact.Scope = append(icyfact.Scope, icy)
	icyfact.Card = append(icyfact.Card, icy.Numvalues)

	holmes := myNet.GetNode("Holmes")

	holmesfact := new(bn.FactorV2)
	holmesfact.CPT = holmes.CPT
	holmesfact.Scope = append(holmesfact.Scope, holmes)
	holmesfact.Scope = append(holmesfact.Scope, holmes.Parents...)
	holmesfact.Card = append(holmesfact.Card, holmes.Numvalues)
	for _, p := range holmes.Parents {
		holmesfact.Card = append(holmesfact.Card, p.Numvalues)
	}

	watson := myNet.GetNode("Watson")
	watsonfact := new(bn.FactorV2)
	watsonfact.CPT = holmes.CPT
	watsonfact.Scope = append(watsonfact.Scope, watson)
	watsonfact.Scope = append(watsonfact.Scope, watson.Parents...)
	watsonfact.Card = append(watsonfact.Card, watson.Numvalues)
	for _, p := range watson.Parents {
		watsonfact.Card = append(watsonfact.Card, p.Numvalues)
	}

	fmt.Println("F V2 Icy", icyfact)
	fmt.Println("F V2 Holmes", holmesfact)
	fmt.Println("F V2 Watson", watsonfact)

	ih := bn.MultiplyFactor(icyfact, holmesfact, true)
	fmt.Println("F V2 Icy", icyfact)
	fmt.Println("F V2 Holmes", holmesfact)
	fmt.Println("Icy-Holmes", ih)

	ihClique := bn.NewClique(ih, "ih")

	iw := bn.MultiplyFactor(icyfact, watsonfact, true)
	fmt.Println("F V2 Icy", icyfact)
	fmt.Println("F V2 Watson", watsonfact)
	fmt.Println("Icy-Watson", iw)

	iwClique := bn.NewClique(iw, "iw")

	jt := new(bn.JunctionTree)
	jt.AddCliques(ihClique, iwClique)
	isep := bn.NewSeparator(icyfact, ihClique, iwClique)
	jt.AddSeparators(isep)
	fmt.Println(isep)
	jt.SetRoot()
	jt.Propagate()

	// bn.EnterEvidenceFactor(ih, icy, []float64{0, 1})
	jt.EnterEvidence(icy, []float64{0, 1})
	jt.Propagate()
	for _, c := range jt.Cliques {
		fmt.Println(c.Name, c.Table)
	}
	// ih.Marginalize(false, icy)
	// iw.Marginalize(false, watson)
	// h := bn.DivideFactor(&ih, holmesfact)
	// fmt.Println(h)
	// holmesEV := bn.EnterEvidence(&ih, holmes, false)
	// fmt.Println("EVID", holmesEV)

	// res := bn.MultiplyFactor(&ih, holmesEV)
	// fmt.Println(ih)
	// fmt.Println(res)
}
