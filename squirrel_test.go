package main

import (
	"testing"
)

func TestPhi(t *testing.T) {
	input := Counts{
		n00: 1,
		n11: 1,
		n01: 1,
		n10: 1,
	}

	expected := 0.0
	result := phi(input)

	if expected != result {
		t.Errorf("Expected 0.0 but hot %v", result)
	}

	input = Counts{
		n00: 20,
		n11: 20,
		n01: 0,
		n10: 0,
	}

	expected = 1.0
	result = phi(input)

	if expected != result {
		t.Errorf("Expected 0.0 but hot %v", result)
	}

}
