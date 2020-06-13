package examples

import (
	"fmt"
	"testing"
	"time"

	bn "github.com/Petrus97/bif-parser/bayesnet"
	"github.com/Petrus97/bif-parser/parser"
)

func TestIcyRoads(t *testing.T) {
	start := time.Now()
	myNet := parser.ReadBIF("../data/icy_roads.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()

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
	w := ih.Marginalize(true, icy)
	i := ih.Marginalize(true, holmes)
	h := iw.Marginalize(true, icy)
	icy.Prob.Prob = i.CPT
	if icy.Prob.Prob[0] != 0.7 {
		t.Error("test failed", icy.Prob.Prob)
	}
	holmes.Prob.Prob = h.CPT
	if holmes.Prob.Prob[0] != 0.59 {
		t.Error("test failed", holmes.Prob.Prob)
	}
	watson.Prob.Prob = w.CPT
	if watson.Prob.Prob[0] != 0.59 {
		t.Error("test failed", watson.Prob.Prob)
	}
	myNet.ListNodes()
	fmt.Println("####################### END TEST ############################")
}

func TestIcyRoadsWithEvidence(t *testing.T) {
	start := time.Now()
	myNet := parser.ReadBIF("../data/icy_roads.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()

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

	jt.SetRoot()
	jt.EnterEvidence(icy, []float64{1, 0}) // Icy = true
	jt.Propagate()
	for _, c := range jt.Cliques {
		fmt.Println(c.Name, c.Table)
	}
	h := ih.Marginalize(true, icy)
	i := ih.Marginalize(true, holmes)
	w := iw.Marginalize(true, icy)
	icy.Prob.Prob = i.CPT
	if icy.Prob.Prob[0] != 1 {
		t.Error("test failed")
	}
	holmes.Prob.Prob = h.CPT
	if holmes.Prob.Prob[0] != 0.8 {
		t.Error("test failed")
	}
	watson.Prob.Prob = w.CPT
	if watson.Prob.Prob[0] != 0.8 {
		t.Error("test failed")
	}
	myNet.ListNodes()
	fmt.Println("####################### END TEST ############################")
}
