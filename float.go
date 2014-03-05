package stl

// This file contains generic numerical functions for floating point processing

import (
	"math"
)

// Just for convenience.
const Pi = math.Pi

// Just for convenience.
const TwoPi = math.Pi * 2

// Just for convenience.
const HalfPi = math.Pi * 0.5

// Just for convenience.
const QuarterPi = math.Pi * 0.25

const tolerance = float32(0.0005)

// Returns true, if a and b are equal allowing for numerical error tol.
func almostEqual32(a, b, tol float32) bool {
	return math.Abs(float64(a-b)) <= float64(tol) 
}

// Returns true, if a and b are equal allowing for numerical error tol.
func almostEqual64(a, b, tol float64) bool {
  return math.Abs(a-b) <= tol
}

func min(a, b float32) float32 {
	if a < b {
		return a
	} else {
		return b
	}
}

func min4(a, b, c, d float32) float32 {
	return min(min(a, b), min(c, d))
}

func max(a, b float32) float32 {
	if a > b {
		return a
	} else {
		return b
	}
}

func max4(a, b, c, d float32) float32 {
	return max(max(a, b), max(c, d))
}
