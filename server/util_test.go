package main

import "testing"

func TestParsePosition(t *testing.T) {
	inputString := "1.21235,2.1241"
	outputExpected := Position{1.21235, 2.1241}
	outputReal, err := ParsePosition(inputString)
	if err != nil {
		t.Errorf("Conversion resulted in unexpected error", err.Error())
	}
	if outputReal != outputExpected {
		t.Errorf("Should be %v but is %v", outputExpected, outputReal)
	}
}
