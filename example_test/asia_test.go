package examples

import (
	"fmt"
	"testing"
	"time"

	bn "github.com/Petrus97/bif-parser/bayesnet"
	"github.com/Petrus97/bif-parser/parser"
)

func TestAsia(t *testing.T) {
	start := time.Now()
	myNet := parser.ReadBIF("../data/asia.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
	// Asia
	asia := myNet.GetNode("asia")
	tub := myNet.GetNode("tub")
	smoke := myNet.GetNode("smoke")
	lung := myNet.GetNode("lung")
	bronc := myNet.GetNode("bronc")
	either := myNet.GetNode("either")
	xray := myNet.GetNode("xray")
	dysp := myNet.GetNode("dysp")

	asiafact := bn.CreateFactorV2(asia)
	tubfact := bn.CreateFactorV2(tub)
	smokefact := bn.CreateFactorV2(smoke)
	lungfact := bn.CreateFactorV2(lung)
	broncfact := bn.CreateFactorV2(bronc)
	eitherfact := bn.CreateFactorV2(either)
	xrayfact := bn.CreateFactorV2(xray)
	dyspfact := bn.CreateFactorV2(dysp)

	fmt.Println(asiafact)
	fmt.Println(tubfact)
	fmt.Println(smokefact)
	fmt.Println(lungfact)
	fmt.Println(broncfact)
	fmt.Println(eitherfact)
	fmt.Println(xrayfact)
	fmt.Println(dyspfact)

	// // P(W, R) = P(W|R)P(R)
	// wr := bn.MultiplyFactor(rainfact, watsonfact, true)
	// fmt.Println("WR", wr)
	// wrClique := bn.NewClique(wr, "wr")

	// // P(H, R, S) = P(H|R,S)P(R,S) = P(H|R,S)P(R)P(S)
	// rs := bn.MultiplyFactor(sprinklerfact, rainfact, true)
	// fmt.Println("RS", rs)
	// hrs := bn.MultiplyFactor(rs, holmesfact, true)
	// fmt.Println("HRS", hrs)

	// hrsClique := bn.NewClique(hrs, "hrs")

	// jt := new(bn.JunctionTree)
	// jt.AddCliques(hrsClique, wrClique)
	// rsep := bn.NewSeparator(watsonfact, wrClique, hrsClique)
	// jt.AddSeparators(rsep)
	// jt.SetRoot()
	// jt.Propagate()

	// jt.EnterEvidence(holmes, []float64{1, 0})
	// jt.Propagate()
	// for _, c := range jt.Cliques {
	// 	fmt.Println(c.Name, c.Table)
	// }

}
