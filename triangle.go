package stl

// This file defines the Triangle data type, the building block for Solid

// Data type for single triangles used in Solid.Triangles. The vertices
// have to be ordered counter-clockwise when looking at their outside surface.
// The vector Normal is orthogonal to the triangle, pointing outside, and
// has length 1. This is redundant but included in the STL format in order to
// avoid recalculation.
type Triangle struct {
	// Normal vector of triangle, should be normalized...
	Normal Vec3

	// Vertices of triangle in right hand order.
	// I.e. from the front the triangle's vertices are ordered counterclockwise
	// and the normal vector is orthogonal to the front pointing outside.
	Vertices [3]Vec3

	// 16 bits of attributes. Not available in ASCII format. Could be used
	// for color selection, texture selection, refraction etc. Some tools ignore
	// this field completely, always writing 0 on export.
	Attributes uint16
}

// Calculate the normal vector using the right hand rule
func (t *Triangle) calculateNormal() Vec3 {
  // The normal is calculated by normalizing the result of
	// (V0-V2) x (V1-V2)
	return t.Vertices[0].diff(t.Vertices[2]).
		cross(t.Vertices[1].diff(t.Vertices[2])).
		unitVec()
}

// Recalculate the redundant normal vector using the right hand rule
func (t *Triangle) recalculateNormal() {
	t.Normal = t.calculateNormal()
}

// Applies a 4x4 transformation matrix to every vertex
// and recalculates the normal
func (t *Triangle) transform(transformationMatrix *Mat4) {
	t.transformNR(transformationMatrix)
	t.recalculateNormal()
}

// Applies a 4x4 transformation matrix to every vertex
// without recalculating the normal afterwards
func (t *Triangle) transformNR(transformationMatrix *Mat4) {
	t.Vertices[0] = transformationMatrix.MultVec3(t.Vertices[0])
	t.Vertices[1] = transformationMatrix.MultVec3(t.Vertices[1])
	t.Vertices[2] = transformationMatrix.MultVec3(t.Vertices[2])
}

// Returns true if at least two vertices are exactly equal, meaning
// this is a line, or even a dot.
func (t *Triangle) hasEqualVertices() bool {
	return t.Vertices[0] == t.Vertices[1] ||
		t.Vertices[0] == t.Vertices[2] ||
		t.Vertices[1] == t.Vertices[2]
}

// Checks if normal matches vertices using right hand rule, with
// numerical tolerance for angle between them given by tol in radians.
func (t *Triangle) checkNormal(tol float64) bool {
  calculatedNormal := t.calculateNormal()
  return t.Normal.angle(calculatedNormal) < tol
}
