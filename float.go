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

// almostEqual returns true, if a and b are equal allowing for numerical error tol.
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
	}
	return b
}

func min4(a, b, c, d float32) float32 {
	return min(min(a, b), min(c, d))
}

func max(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

func max4(a, b, c, d float32) float32 {
	return max(max(a, b), max(c, d))
}
