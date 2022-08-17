package stl

// This file contains the 3D vector data type that is used for the triangles

import (
	"math"
)

// Vec3 represents a 3D vector, used in Triangle for normal vector and vertices.
type Vec3 [3]float32

// vec3Zero is the zero vector
var vec3Zero = Vec3{0, 0, 0}

// len returns the Euclidean length of a vector.
func (vec Vec3) len() float64 {
	return math.Sqrt(float64(vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]))
}

// UnitVec3 returns vec multiplied by 1/vec.Len(), so its new length is 1. If the vector is empty, it is returned as such.
func (vec Vec3) UnitVec3() Vec3 {
	l := vec.len()
	if l == 0 {
		return vec
	}

	return Vec3{float32(float64(vec[0]) / l), float32(float64(vec[1]) / l), float32(float64(vec[2]) / l)}
}

// MultScalar multiplies vec by scalar.
func (vec Vec3) MultScalar(scalar float64) Vec3 {
	return Vec3{float32(float64(vec[0]) * scalar), float32(float64(vec[1]) * scalar), float32(float64(vec[2]) * scalar)}
}

// AlmostEqual returns true if vec and o are equal allowing for numerical error tol.
func (vec Vec3) AlmostEqual(o Vec3, tol float32) bool {
	return almostEqual32(vec[0], o[0], tol) && almostEqual32(vec[1], o[1], tol) && almostEqual32(vec[2], o[2], tol)
}

// Add returns the sum of vectors vec and o.
func (vec Vec3) Add(o Vec3) Vec3 {
	return Vec3{
		vec[0] + o[0],
		vec[1] + o[1],
		vec[2] + o[2],
	}
}

// Diff returns the difference vec - o.
func (vec Vec3) Diff(o Vec3) Vec3 {
	return Vec3{
		vec[0] - o[0],
		vec[1] - o[1],
		vec[2] - o[2],
	}
}

// Cross returns the vector Cross product vec x o.
func (vec Vec3) Cross(o Vec3) Vec3 {
	return Vec3{
		vec[1]*o[2] - vec[2]*o[1],
		vec[2]*o[0] - vec[0]*o[2],
		vec[0]*o[1] - vec[1]*o[0],
	}
}

// Dot returns the Dot product between vec and o.
func (vec Vec3) Dot(o Vec3) float64 {
	return float64(vec[0])*float64(o[0]) +
		float64(vec[1])*float64(o[1]) +
		float64(vec[2])*float64(o[2])
}

// Angle between vec and o in radians, without sign, between 0 and Pi.
// If vec or o is the origin, this returns 0.
func (vec Vec3) Angle(o Vec3) float64 {
	lenProd := vec.len() * o.len()
	if lenProd == 0 {
		return 0
	}
	cosAngle := vec.Dot(o) / lenProd
	// Numerical correction
	if cosAngle < -1 {
		cosAngle = -1
	} else if cosAngle > 1 {
		cosAngle = 1
	}

	return math.Acos(cosAngle)
}
