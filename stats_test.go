package main

import (
	"math"
	"testing"
)

func TestCombineAvg(t *testing.T) {
	testAvg := combineAvg(.6, 16, .2, 24)
	testEqualFloats(t, testAvg, 0.36)
}

func TestCombineSD(t *testing.T) {
	testSD := combineSD(1.479019946, 5.75, 4, 1.124858268, 5.857142857, 7)
	testEqualFloats(t, testSD, 1.266217)
}

func TestSD(t *testing.T) {
	testSD := sd(390, 5.818181818, 11)
	testEqualFloats(t, testSD, 1.266217)
}

func testEqualFloats(t *testing.T, actual float64, want float64) {
	if math.Abs(want - actual) > .001 {
		t.Errorf("Wanted %f got %f", want, actual) 
	}
}

