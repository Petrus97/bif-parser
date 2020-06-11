package examples

import (
	"fmt"
	"testing"
	"time"

	bn "github.com/Petrus97/bif-parser/bayesnet"
	"github.com/Petrus97/bif-parser/parser"
)

func TestWetGrass(t *testing.T) {
	start := time.Now()
	myNet := parser.ReadBIF("../data/wet_grass.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
	// WET GRASS
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

	// P(W, R) = P(W|R)P(R)
	wr := bn.MultiplyFactor(rainfact, watsonfact, true)
	fmt.Println("WR", wr)
	wrClique := bn.NewClique(wr, "wr")

	// P(H, R, S) = P(H|R,S)P(R,S) = P(H|R,S)P(R)P(S)
	rs := bn.MultiplyFactor(sprinklerfact, rainfact, true)
	fmt.Println("RS", rs)
	hrs := bn.MultiplyFactor(rs, holmesfact, true)
	fmt.Println("HRS", hrs)

	hrsClique := bn.NewClique(hrs, "hrs")

	jt := new(bn.JunctionTree)
	jt.AddCliques(hrsClique, wrClique)
	rsep := bn.NewSeparator(watsonfact, wrClique, hrsClique)
	jt.AddSeparators(rsep)
	jt.SetRoot()
	jt.Propagate()

	jt.EnterEvidence(holmes, []float64{1, 0})
	jt.Propagate()
	for _, c := range jt.Cliques {
		fmt.Println(c.Name, c.Table)
	}
	w := wrClique.Table.Marginalize(true, watson)
	fmt.Println("w", w)

}
