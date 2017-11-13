package stl

// This file contains a 4x4 matrix implementation for 3D transformations

// Mat4 represents a 4x4 Matrix of float64 used for 3D transformations.
// The 4th column can be used for moving the solid on the axes.
// Accessing matrix elements goes like this:
//    matrix[row][column]
type Mat4 [4]Vec4

// Vec4 is used to construct Mat4
type Vec4 [4]float64

// MultMat4 Multiplies m with o and write the result into r.
func (m *Mat4) MultMat4(o *Mat4, r *Mat4) {
	// I tried Strassen here, but this is faster - it is simply the number of ops
	// that counts, not just the number of multiplications.
	// Manual loop-unrolling
	r[0][0] = m[0][0]*o[0][0] + m[0][1]*o[1][0] + m[0][2]*o[2][0] + m[0][3]*o[3][0]
	r[0][1] = m[0][0]*o[0][1] + m[0][1]*o[1][1] + m[0][2]*o[2][1] + m[0][3]*o[3][1]
	r[0][2] = m[0][0]*o[0][2] + m[0][1]*o[1][2] + m[0][2]*o[2][2] + m[0][3]*o[3][2]
	r[0][3] = m[0][0]*o[0][3] + m[0][1]*o[1][3] + m[0][2]*o[2][3] + m[0][3]*o[3][3]
	r[1][0] = m[1][0]*o[0][0] + m[1][1]*o[1][0] + m[1][2]*o[2][0] + m[1][3]*o[3][0]
	r[1][1] = m[1][0]*o[0][1] + m[1][1]*o[1][1] + m[1][2]*o[2][1] + m[1][3]*o[3][1]
	r[1][2] = m[1][0]*o[0][2] + m[1][1]*o[1][2] + m[1][2]*o[2][2] + m[1][3]*o[3][2]
	r[1][3] = m[1][0]*o[0][3] + m[1][1]*o[1][3] + m[1][2]*o[2][3] + m[1][3]*o[3][3]
	r[2][0] = m[2][0]*o[0][0] + m[2][1]*o[1][0] + m[2][2]*o[2][0] + m[2][3]*o[3][0]
	r[2][1] = m[2][0]*o[0][1] + m[2][1]*o[1][1] + m[2][2]*o[2][1] + m[2][3]*o[3][1]
	r[2][2] = m[2][0]*o[0][2] + m[2][1]*o[1][2] + m[2][2]*o[2][2] + m[2][3]*o[3][2]
	r[2][3] = m[2][0]*o[0][3] + m[2][1]*o[1][3] + m[2][2]*o[2][3] + m[2][3]*o[3][3]
	r[3][0] = m[3][0]*o[0][0] + m[3][1]*o[1][0] + m[3][2]*o[2][0] + m[3][3]*o[3][0]
	r[3][1] = m[3][0]*o[0][1] + m[3][1]*o[1][1] + m[3][2]*o[2][1] + m[3][3]*o[3][1]
	r[3][2] = m[3][0]*o[0][2] + m[3][1]*o[1][2] + m[3][2]*o[2][2] + m[3][3]*o[3][2]
	r[3][3] = m[3][0]*o[0][3] + m[3][1]*o[1][3] + m[3][2]*o[2][3] + m[3][3]*o[3][3]
}

// Mat4Identity is the identity matrix
var Mat4Identity = Mat4{
	Vec4{1, 0, 0, 0},
	Vec4{0, 1, 0, 0},
	Vec4{0, 0, 1, 0},
	Vec4{0, 0, 0, 1},
}

// MultVec3 multiplies m with v, where v[3] is assumed to be 1, and the 4th result value is
// not calculated, as is usual in 3D transformations.
func (m *Mat4) MultVec3(v Vec3) Vec3 {
	var result Vec3
	result[0] = float32(m[0][0]*float64(v[0]) + m[0][1]*float64(v[1]) + m[0][2]*float64(v[2]) + m[0][3])
	result[1] = float32(m[1][0]*float64(v[0]) + m[1][1]*float64(v[1]) + m[1][2]*float64(v[2]) + m[1][3])
	result[2] = float32(m[2][0]*float64(v[0]) + m[2][1]*float64(v[1]) + m[2][2]*float64(v[2]) + m[2][3])
	return result
}
