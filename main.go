package main

import (
	"fmt"
	"time"

	"github.com/Petrus97/bif-parser/parser"
)

func main() {
	fmt.Println("Hello")
	start := time.Now()
	myNet := parser.ReadBIF("data/survey.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
	myNet.ListNodes()
}
