package stl

// This file contains tests for the Vec3 data type

import (
	"testing"
)

func TestVec3Angle(t *testing.T) {
	v := Vec3{1, 0, 0}
	tol := 0.00005
	testV := []Vec3{
		{0, 1, 0},
		{0, -1, 0},
		{-1, 0, 0},
		{-1, 1, 0},
		{-1, -1, 0}}
	expected := []float64{
		HalfPi,
		HalfPi,
		Pi,
		HalfPi + QuarterPi,
		HalfPi + QuarterPi}
	for i, tv := range testV {
		r := v.Angle(tv)
		if !almostEqual64(expected[i], r, tol) {
			t.Errorf("Angle(%v, %v) = %g Pi, expected %g Pi", v, tv, r/Pi, expected[i]/Pi)
		}
	}
}
