package stl

// Writer processes an STL solid as a stream.
type Writer interface {
	// SetName sets the solid's name
	SetName(name string)

	// SetBinaryHeader explicitly sets a binary header used in binary STL
	SetBinaryHeader(header []byte)

	// SetASCII sets the IsAscii flag that indicates whether the solid was read from ASCII STL
	SetASCII(isASCII bool)

	// SetTriangleCount can optionally be used to set the triangle count if it is known, so
	// the underlying implementation can use it to allocate a data structure, or similar.
	SetTriangleCount(n uint32)

	// AppendTriangle adds a triangle to the solid
	AppendTriangle(t Triangle)
}
