package stl

// This file provides the Solid data type that is a memory representation
// of an STL file

// Solid is a 3D model made out of triangles, called solid in STL, representing
// an STL file
type Solid struct {
	// only used in binary format
	BinaryHeader []byte

	Name string

	Triangles []Triangle

	// true, if this Solid was read from an ASCII file, and false, if read
	// from a binary file. Also used to determine the format when writing
	// to a file.
	IsAscii bool
}

// SolidMeasure is used to store the result of Solid.Measure()
type SolidMeasure struct {
	// Minimum values for axes
	Min Vec3

	// Maximum values for axes
	Max Vec3

	// Max - Min
	Len Vec3
}

// Measure the dimensions of a solid in its own units
func (solid *Solid) Measure() SolidMeasure {
	var measure SolidMeasure

	if len(solid.Triangles) == 0 {
		return measure
	}

	// initialize with real values
	a := &solid.Triangles[0].Vertices[0]
	measure.Min = *a
	measure.Max = *a

	for _, triangle := range solid.Triangles {
		for d := 0; d < 3; d++ {
			measure.Min[d] = min4(measure.Min[d], triangle.Vertices[0][d], triangle.Vertices[1][d], triangle.Vertices[2][d])
			measure.Max[d] = max4(measure.Max[d], triangle.Vertices[0][d], triangle.Vertices[1][d], triangle.Vertices[2][d])
		}
	}

	measure.Len = measure.Max.diff(measure.Min)

	return measure
}

// Transform applies a 4x4 transformation matrix to every vertex
// and recalculates the normal for every triangle
func (solid *Solid) Transform(transformationMatrix *Mat4) {
	l := len(solid.Triangles)
	for i := 0; i < l; i++ {
		// Tried go-routines here. Was slower even with large solids.
		solid.Triangles[i].transform(transformationMatrix)
	}
}

// TransformNR applies a 4x4 transformation matrix to every vertex
// and does not recalculate the normal vector for every triangle.
// This could be used to speed things up when multiple transformations
// are applied successively to a solid, and the transformation matrix
// is not calculated beforehand. Before writing this solid to disk then,
// RecalculateNormals() should be called.
func (solid *Solid) TransformNR(transformationMatrix *Mat4) {
	l := len(solid.Triangles)
	for i := 0; i < l; i++ {
		// Tried go-routines here. Was slower even with large solids.
		solid.Triangles[i].transformNR(transformationMatrix)
	}
}

// RecalculateNormals recalculates all triangle normal vectors from the vertices.
// Can be used after multiple transformations using the TransformNR method
// that does not recalculate the normal vectors.
func (solid *Solid) RecalculateNormals() {
	for i := 0; i < len(solid.Triangles); i++ {
		solid.Triangles[i].recalculateNormal()
	}
}

// Scale all vertex coordinates by scalar factor
func (solid *Solid) Scale(factor float64) {
	for i := 0; i < len(solid.Triangles); i++ {
		t := &solid.Triangles[i]
		for v := 0; v < 3; v++ {
			for d := 0; d < 3; d++ {
				t.Vertices[v][d] = float32(factor * float64(t.Vertices[v][d]))
			}
		}
	}
}

// Stretch scales all vertex coordinates by different factors per axis
func (solid *Solid) Stretch(vec Vec3) {
	for i := 0; i < len(solid.Triangles); i++ {
		t := &solid.Triangles[i]
		for v := 0; v < 3; v++ {
			for d := 0; d < 3; d++ {
				t.Vertices[v][d] = float32(float64(vec[d]) * float64(t.Vertices[v][d]))
			}
		}
		t.recalculateNormal()
	}
}

// ScaleLinearDowntoSizeBox works like this: if the solid does not fit into size box
// defined by sizeBox, it is scaled down accordingly. It is not scaled up, if it is
// smaller than sizeBox. All sizes have to be > 0.
func (solid *Solid) ScaleLinearDowntoSizeBox(sizeBox Vec3) {
	if sizeBox[0] <= 0 || sizeBox[1] <= 0 || sizeBox[2] <= 0 {
		panic("Not all values in sizeBox are > 0!")
	}
	measure := solid.Measure()
	factor := float64(1)
	for d := 0; d < 3; d++ {
		if measure.Len[d] > sizeBox[d] {
			scale := float64(sizeBox[d]) / float64(measure.Len[d])
			if scale < factor {
				factor = scale
			}
		}
	}
	if factor != float64(1) {
		solid.Scale(factor)
	}
}

// Translate (i.e. move) the solid by vec
func (solid *Solid) Translate(vec Vec3) {
	for i := 0; i < len(solid.Triangles); i++ {
		t := &solid.Triangles[i]
		for v := 0; v < 3; v++ {
			for d := 0; d < 3; d++ {
				t.Vertices[v][d] += vec[d]
			}
		}
	}
}

// IsInPositive is true if every vertex in this solid is within the positive octant, i.e.
// all coordinate values are positive or 0.
func (solid *Solid) IsInPositive() bool {
	measure := solid.Measure()
	for dim := 0; dim < 3; dim++ {
		if measure.Min[dim] < 0 {
			return false
		}
	}
	return true
}

// MoveToPositive moves the solid into the positive octant if necessary, as prescribed by
// the original STL format spec. Some applications tolerate negative coordinates. This also
// makes sense, as the origin is a perfect reference point for rotations.
func (solid *Solid) MoveToPositive() {
	measure := solid.Measure()
	var translationVector Vec3
	for dim := 0; dim < 3; dim++ {
		if measure.Min[dim] < 0 {
			translationVector[dim] = -measure.Min[dim] // move smallest value to 0
		}
	}
	// only apply vector if non-zero
	if translationVector != vec3Zero {
		solid.Translate(translationVector)
	}
}

// scaleDim scales only dimension dim of solid by scalar factor.
func (solid *Solid) scaleDim(factor float64, dim int) {
	for i := 0; i < len(solid.Triangles); i++ {
		for v := 0; v < 3; v++ {
			solid.Triangles[i].Vertices[v][dim] = float32(factor * float64(solid.Triangles[i].Vertices[v][dim]))
		}
	}
}

// Rotate the solid by angle radians around a rotation axis defined
// by a point pos on the axis and a direction vector dir. This
// example would rotate the solid by 90 degree around the z-axis:
//    stl.Rotate(stl.Vec3{0,0,0}, stl.Vec3{0,0,1}, stl.HalfPi)
func (solid *Solid) Rotate(pos, dir Vec3, angle float64) {
	var rotationMatrix Mat4
	RotationMatrix(pos, dir, angle, &rotationMatrix)
	// Apply rotationMatrix and recalculate normal vectors
	solid.Transform(&rotationMatrix)
}

// TriangleErrors represent the errors found in a single triangle.
type TriangleErrors struct {
	// HasEqualVertices is true if some vertices are identical, meaning we are having
	// a line, or even a point, as opposed to a triangle.
	HasEqualVertices bool

	// NormalDoesNotMatch istrue if the normal vector does not match a normal calculated from the
	// vertices in the right hand order, even allowing for an angular difference
	// of < 90 degree.
	NormalDoesNotMatch bool

	// EdgeErrors by edge. The edge is indexed by it's first vertex, i.e.
	//    0: V0 -> V1
	//    1: V1 -> V2
	//    2: V2 -> V0
	// If the edge has no error its value is nil.
	EdgeErrors [3]*EdgeError
}

// edge is a convenience accessor that allocates an EdgeError for edge e if
// it is not already present. e is the index of the edge, being equal
// to the index of the first vertex.
func (te *TriangleErrors) edge(e int) *EdgeError {
	if te.EdgeErrors[e] == nil {
		te.EdgeErrors[e] = new(EdgeError)
	}
	return te.EdgeErrors[e]
}

// EdgeError describes the errors found for a single edge within a triangle using
// Solid.Validate().
type EdgeError struct {
	// SameEdgeTriangles are indexes in Solid.Triangles of triangles that contain exactly the same edge.
	SameEdgeTriangles []int

	// CounterEdgeTriangles are indexes in Solid.Triangles of triangles that contain the edge in the
	// opposite direction. If there is exactly one other triangle, this is no
	// error.
	CounterEdgeTriangles []int
}

// IsUsedInOtherTriangles is true if this edge is also used in another triangle, meaning that
// there is probably something wrong with this or the other triangle's
// orientation.
func (eer *EdgeError) IsUsedInOtherTriangles() bool {
	return len(eer.SameEdgeTriangles) != 0
}

// HasMultipleCounterEdges is true if there is more than one other triangle
// with this edge in the opposite direction
func (eer *EdgeError) HasMultipleCounterEdges() bool {
	return len(eer.CounterEdgeTriangles) > 1
}

// HasNoCounterEdge is true if there is no other triangle with this edge in the opposite
// direction, meaning that there is no neighboring triangle
func (eer *EdgeError) HasNoCounterEdge() bool {
	return len(eer.CounterEdgeTriangles) == 0
}

// For every edge described by two points this data structure stores
// the set of indices of triangles containing this edge.
type edgeLookup struct {
	edgeToTriangles map[[2]Vec3](map[int]bool)
}

func newEdgeLookup() *edgeLookup {
	var l edgeLookup
	l.edgeToTriangles = make(map[[2]Vec3](map[int]bool))
	return &l
}

// Get the indices to Solid.Triangles of triangles containing
// the edge v -> w, excluding a specific one denoted by i.
func (l *edgeLookup) OtherTrianglesWithEdge(v, w Vec3, i int) []int {
	triangleSet, found := l.edgeToTriangles[[2]Vec3{v, w}]
	if !found {
		return nil
	}
	r := make([]int, 0, len(triangleSet))
	for t := range triangleSet {
		if t != i {
			r = append(r, t)
		}
	}
	return r
}

// Put triangleIndex into the set of indices of triangles containing
// the edge v -> w.
func (l *edgeLookup) InsertEdge(triangleIndex int, v, w Vec3) {
	key := [2]Vec3{v, w}
	triangleSet, found := l.edgeToTriangles[key]
	if !found {
		l.edgeToTriangles[key] = make(map[int]bool)
		triangleSet = l.edgeToTriangles[key]
	}
	triangleSet[triangleIndex] = true
}

// triangleErrorsMap represents errors by triangle index
type triangleErrorsMap map[int]*TriangleErrors

// item is a convenient map accessor that creates the entry on first use.
func (m triangleErrorsMap) item(triangleIdx int) *TriangleErrors {
	if _, found := m[triangleIdx]; !found {
		m[triangleIdx] = new(TriangleErrors)
	}
	return m[triangleIdx]
}

const normalAngleTolerance = HalfPi

// Validate looks for triangles that are really lines or dots, and for edges that
// violate the vertex-to-vertex rule. Returns a map of errors by triangle
// index that could be used to print out an error report.
func (solid *Solid) Validate() map[int]*TriangleErrors {
	// Build up lookup from edge to triangle
	e := newEdgeLookup()
	for i := range solid.Triangles {
		t := &solid.Triangles[i]
		for vertex1 := 0; vertex1 < 3; vertex1++ {
			vertex2 := (vertex1 + 1) % 3
			e.InsertEdge(i, t.Vertices[vertex1], t.Vertices[vertex2])
		}
	}

	// Now look for every edge that there is exactly one
	// "counter-edge" in the opposite direction for another triangle,
	// and that the same edge is not used by another triangle.
	triangleErrors := make(triangleErrorsMap)

	for i := range solid.Triangles {
		t := &solid.Triangles[i]
		// check for equal vertices
		if t.hasEqualVertices() {
			triangleErrors.item(i).HasEqualVertices = true
		}

		// check if normal matches vertices
		if !t.checkNormal(normalAngleTolerance) {
			triangleErrors.item(i).NormalDoesNotMatch = true
		}

		// loop through edges, vertex1 is also the index of the edge
		for vertex1 := 0; vertex1 < 3; vertex1++ {
			vertex2 := (vertex1 + 1) % 3

			// look for other triangles with same edge
			sameEdgeTriangles := e.OtherTrianglesWithEdge(
				t.Vertices[vertex1], t.Vertices[vertex2], i)
			if len(sameEdgeTriangles) > 0 {
				// used by other triangles
				triangleErrors.item(i).edge(vertex1).SameEdgeTriangles = sameEdgeTriangles
			}

			// Look for other triangles with edge in opposite direction.
			// If this same triangle had the edge in the opposite direction, we
			// would already have the HasEqualVertices error, so we do not need
			// to take care of this here.
			counterEdgeTriangles := e.OtherTrianglesWithEdge(
				t.Vertices[vertex2], t.Vertices[vertex1], i)
			if len(counterEdgeTriangles) != 1 {
				// 0 or multiple "counter-edges"
				triangleErrors.item(i).edge(vertex1).CounterEdgeTriangles = counterEdgeTriangles
			}
		}
	}

	return triangleErrors
}
