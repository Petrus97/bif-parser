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

	// bn.EnterEvidenceFactor(ih, icy, []float64{0, 1})
	jt.EnterEvidence(icy, []float64{1, 0})
	jt.Propagate()
	for _, c := range jt.Cliques {
		fmt.Println(c.Name, c.Table)
	}
	i := ih.Marginalize(true, icy)
	h := ih.Marginalize(true, holmes)
	w := iw.Marginalize(true, watson)
	icy.Prob.Prob = i.CPT
	holmes.Prob.Prob = h.CPT
	watson.Prob.Prob = w.CPT
	myNet.ListNodes()
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
