package stl

// This file contains generic numerical functions for floating point processing

import (
	"math"
)

// Pi is just math.Pi
const Pi = math.Pi

// TwoPi is math.Pi * 2
const TwoPi = math.Pi * 2

// HalfPi is math.Pi * 0.5
const HalfPi = math.Pi * 0.5

// QuarterPi is math.Pi * 0.25
const QuarterPi = math.Pi * 0.25

// AlmostEqual returns true, if a and b are equal allowing for numerical error tol.
func almostEqual32(a, b, tol float32) bool {
	return math.Abs(float64(a-b)) <= float64(tol)
}

// Returns true, if a and b are equal allowing for numerical error tol.
func almostEqual64(a, b, tol float64) bool {
	return math.Abs(a-b) <= tol
}
