package stl

// This file defines functions to write a Solid into the STL binary format.

import (
	"encoding/binary"
	"io"
	"math"
)

// Write solid in binary STL into an io.Writer.
// Does not check whether len(solid.Triangles) fits into uint32.
func writeSolidBinary(w io.Writer, solid *Solid) error {
	headerBuf := make([]byte, binaryHeaderSize)
	if solid.BinaryHeader == nil {
		// use name if no binary header set
		copy(headerBuf, solid.Name)
	} else {
		copy(headerBuf, solid.BinaryHeader)
	}
	// Write triangle count
	binary.LittleEndian.PutUint32(headerBuf[80:84], uint32(len(solid.Triangles)))
	_, errHeader := w.Write(headerBuf)
	if errHeader != nil {
		return errHeader
	}

	// Write each triangle
	for _, t := range solid.Triangles {
		tErr := writeTriangleBinary(w, &t)
		if tErr != nil {
			return tErr
		}
	}

	return nil
}

func writeTriangleBinary(w io.Writer, t *Triangle) error {
	buf := make([]byte, 50)
	offset := 0
	encodePoint(buf, &offset, &t.Normal)
	encodePoint(buf, &offset, &t.Vertices[0])
	encodePoint(buf, &offset, &t.Vertices[1])
	encodePoint(buf, &offset, &t.Vertices[2])
	encodeUint16(buf, &offset, t.Attributes)
	_, err := w.Write(buf)
	return err
}

func encodePoint(buf []byte, offset *int, pt *Vec3) {
	encodeFloat32(buf, offset, pt[0])
	encodeFloat32(buf, offset, pt[1])
	encodeFloat32(buf, offset, pt[2])
}

func encodeFloat32(buf []byte, offset *int, f float32) {
	u32 := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf[*offset:(*offset)+4], u32)
	*offset += 4
}

func encodeUint16(buf []byte, offset *int, u uint16) {
	binary.LittleEndian.PutUint16(buf[*offset:(*offset)+2], u)
	*offset += 2
}
