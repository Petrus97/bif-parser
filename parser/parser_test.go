package parser

import "testing"

func TestParseIcy(t *testing.T) {
	myNet := ReadBIF("../data/icy_roads.bif")
	if len(myNet.Nodes) != 3 {
		t.Error("Error in parsing")
	}
	t.Log("Icy pass")
}
func TestParseGrass(t *testing.T) {
	myNet := ReadBIF("../data/wet_grass.bif")
	if len(myNet.Nodes) != 4 {
		t.Error("Error in parsing")
	}
	t.Log("Wet grass pass")
}

func TestParseMunin(t *testing.T) {
	myNet := ReadBIF("../data/munin.bif")
	if len(myNet.Nodes) != 1041 {
		t.Error("Error in parsing")
	}
	t.Log("Munin pass")
}
