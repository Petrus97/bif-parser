package bayesnet

import (
	"fmt"
	"testing"
)

// Koller matlab example
func TestProductFactor(t *testing.T) {
	node1 := new(Node)
	node1.Name = "1"
	node1.Numvalues = 2
	node1.CPT = []float64{0.11, 0.89}

	node2 := new(Node)
	node2.Name = "2"
	node2.Numvalues = 2
	node2.CPT = []float64{0.59, 0.41, 0.22, 0.78}

	node2.AddParents(node1)

	first := CreateFactorV2(node1)
	second := CreateFactorV2(node2)

	res := MultiplyFactor(first, second, true)
	expectedresults := []float64{0.0649, 0.1958, 0.0451, 0.6942}
	for i := 0; i < len(expectedresults); i++ {
		if expectedresults[i] != res.CPT[i] {
			t.Failed()
		}
	}
	// fmt.Println(res)
}

// Koller fig.9.7 pag 297
func TestMarginalization(t *testing.T) {
	node1 := new(Node)
	node1.Name = "1"
	node1.Numvalues = 3

	node2 := new(Node)
	node2.Name = "2"
	node2.Numvalues = 2

	node3 := new(Node)
	node3.Name = "2"
	node3.Numvalues = 2

	f := new(FactorV2)
	f.CPT = []float64{0.25, 0.35, 0.08, 0.16, 0.05, 0.07, 0., 0., 0.15, 0.21, 0.09, 0.18}
	f.Scope = append(f.Scope, node1)
	f.Scope = append(f.Scope, node2)
	f.Scope = append(f.Scope, node3)
	f.Card = []int{3, 2, 2}

	ret := f.Marginalize(true, node2)
	expectedresults := []float64{0.33, 0.51, 0.05, 0.07, 0.24, 0.39}
	for i := 0; i < len(expectedresults); i++ {
		if expectedresults[i] != f.CPT[i] {
			t.Failed()
		}
	}
	fmt.Println(ret)
}

func TestDivideFactor(t *testing.T) {
	fmt.Println("Test Divide Factor")
	node1 := new(Node)
	node1.Name = "1"
	node1.Numvalues = 3

	node2 := new(Node)
	node2.Name = "2"
	node2.Numvalues = 2

	phi1 := new(FactorV2)
	phi1.CPT = []float64{0.5, 0.2, 0., 0., 0.3, 0.45}
	phi1.Scope = append(phi1.Scope, node1)
	phi1.Scope = append(phi1.Scope, node2)
	phi1.Card = []int{3, 2}

	phi2 := new(FactorV2)
	phi2.CPT = []float64{0.8, 0., 0.6}
	phi2.Scope = append(phi2.Scope, node1)
	phi2.Card = []int{3}

	ret := DivideFactor(phi1, phi2, true)
	fmt.Println("phi1", phi1)
	fmt.Println("phi2", phi2)
	fmt.Println("ret", ret)

}

// func TestDivideFactorV2(t *testing.T) {
// 	fmt.Println("Test Divide Factor")
// 	node1 := new(Node)
// 	node1.Name = "1"
// 	node1.Numvalues = 2

// 	node2 := new(Node)
// 	node2.Name = "2"
// 	node2.Numvalues = 3

// 	node3 := new(Node)
// 	node3.Name = "2"
// 	node3.Numvalues = 2

// 	phi1 := new(FactorV2)
// 	phi1.CPT = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}
// 	phi1.Scope = append(phi1.Scope, node1)
// 	phi1.Scope = append(phi1.Scope, node2)
// 	phi1.Scope = append(phi1.Scope, node3)
// 	phi1.Card = []int{2, 3, 2}

// 	phi2 := new(FactorV2)
// 	phi2.CPT = []float64{0, 1, 2, 3, 4}
// 	phi2.Scope = append(phi2.Scope, node3)
// 	phi2.Scope = append(phi2.Scope, node1)
// 	phi2.Card = []int{2, 2}

// 	ret := DivideFactor(phi1, phi2, true)
// 	fmt.Println("phi1", phi1)
// 	fmt.Println("phi2", phi2)
// 	fmt.Println("ret", ret)

// }
