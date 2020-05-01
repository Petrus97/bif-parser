package main

import (
	"fmt"
	"time"

	parser "github.com/Petrus97/bif-parser/cmd"
	"github.com/Petrus97/bif-parser/morestrings"
)

func main() {
	fmt.Println("Hello")
	fmt.Println(morestrings.ReverseRunes("!oG ,olleH"))
	start := time.Now()
	myNet := parser.ReadBIF("data/munin.bif")
	fmt.Println(time.Now().Sub(start))
	fmt.Println(myNet)
}
