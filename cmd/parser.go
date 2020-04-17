package cmd

import (
	"fmt"
	"io/ioutil"
)

// Probabilities assigned in the model
type Probabilities struct {
	states map[string]float64
}

// Node if net
type Node struct {
	Type      string
	numvalues int
	prob      Probabilities
	CPT       map[int]float64
}

func errorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func readBIF() {
	data, err := ioutil.ReadFile("../data/survey.bif")
	errorCheck(err)
	fmt.Println(string(data))

}
