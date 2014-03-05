package stl

// This file contains the 3D vector data type that is used for the triangles

import (
	"math"
)

// 3D Vector, used in Triangle for normal vector and vertices.
type Vec3 [3]float32

// The zero vector
var vec3Zero = Vec3{0, 0, 0}

// Returns the Euclidean length of a vector.
func (vec Vec3) len() float64 {
	return math.Sqrt(float64(vec[0]*vec[0] + vec[1]*vec[1] + vec[2]*vec[2]))
}

// vec multiplied by 1/vec.Len(), so its new length is 1. If the vector is empty, it is returned as such.
func (vec Vec3) unitVec() Vec3 {
	l := vec.len()
	if l == 0 {
		return vec
	} else {
		return Vec3{float32(float64(vec[0]) / l), float32(float64(vec[1]) / l), float32(float64(vec[2]) / l)}
	}
}

// Multiply by scalar.
func (vec Vec3) multScalar(scalar float64) Vec3 {
	return Vec3{float32(float64(vec[0]) * scalar), float32(float64(vec[1]) * scalar), float32(float64(vec[2]) * scalar)}
}

// Returns true if vec and o are equal allowing for numerical error tol.
func (vec Vec3) almostEqual(o Vec3, tol float32) bool {
	return almostEqual32(vec[0], o[0], tol) && almostEqual32(vec[1], o[1], tol) && almostEqual32(vec[2], o[2], tol)
}

// Return the sum of vectors vec and o.
func (vec Vec3) add(o Vec3) Vec3 {
	return Vec3{
		vec[0] + o[0],
		vec[1] + o[1],
		vec[2] + o[2],
	}
}

// Return the difference vec - o.
func (vec Vec3) diff(o Vec3) Vec3 {
	return Vec3{
		vec[0] - o[0],
		vec[1] - o[1],
		vec[2] - o[2],
	}
}

// Return the vector cross product vec x o.
func (vec Vec3) cross(o Vec3) Vec3 {
	return Vec3{
		vec[1]*o[2] - vec[2]*o[1],
		vec[2]*o[0] - vec[0]*o[2],
		vec[0]*o[1] - vec[1]*o[0],
	}
}

// Return dot product between vec and o.
func (vec Vec3) dot(o Vec3) float64 {
  return float64(vec[0]) * float64(o[0]) +
    float64(vec[1]) * float64(o[1]) +
    float64(vec[2]) * float64(o[2])
}

// Return angle between vec and o in radians, without sign, between 0 and Pi.
// If vec or o is the origin, this returns 0.
func (vec Vec3) angle(o Vec3) float64 {
  lenProd := vec.len() * o.len()
  if lenProd == 0 {
    return 0
  }
  cosAngle := vec.dot(o) / lenProd
  // Numerical correction
  if cosAngle < -1 {
    cosAngle = -1
  } else if cosAngle > 1 {
    cosAngle = 1
  }
  
  return math.Acos(cosAngle)
}