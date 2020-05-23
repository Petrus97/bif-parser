package parser

import "testing"

func TestParseAsia(t *testing.T) {
	myNet := ReadBIF("../data/asia.bif")
	if len(myNet.Nodes) != 8 {
		t.Error("Error in parsing")
	}
}

func TestParseIcy(t *testing.T) {
	myNet := ReadBIF("../data/icy_roads.bif")
	if len(myNet.Nodes) != 3 {
		t.Error("Error in parsing")
	}
}
func TestParseSurvey(t *testing.T) {
	myNet := ReadBIF("../data/survey.bif")
	if len(myNet.Nodes) != 6 {
		t.Error("Error in parsing")
	}
}
func TestParseGrass(t *testing.T) {
	myNet := ReadBIF("../data/wet_grass.bif")
	if len(myNet.Nodes) != 4 {
		t.Error("Error in parsing")
	}
}

func TestParseMunin(t *testing.T) {
	myNet := ReadBIF("../data/munin.bif")
	if len(myNet.Nodes) != 1041 {
		t.Error("Error in parsing")
	}
}
