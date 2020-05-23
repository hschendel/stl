package stl

// This file defines reading functions for the STL binary format.

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

func readAllBinary(r io.Reader) (solid *Solid, err error) {
	header := make([]byte, 84)
	n, readErr := r.Read(header)
	if readErr == io.EOF && n != 84 {
		err = ErrIncompleteBinaryHeader
		return
	} else if readErr != nil {
		err = readErr
		return
	}

	var solidData Solid
	solidData.BinaryHeader = header[0:80]
	solidData.Name = extractASCIIString(solidData.BinaryHeader)
	triangleCount := binary.LittleEndian.Uint32(header[80:84])
	solidData.Triangles = make([]Triangle, triangleCount)

	for i := range solidData.Triangles {
		readErr = readTriangleBinary(r, &solidData.Triangles[i])
		if readErr != nil {
			err = fmt.Errorf("While reading triangle no. %d at byte %d: %s", i, 84+i*50, readErr.Error())
			return
		}
	}

	solid = &solidData
	return
}

func readTriangleBinary(r io.Reader, t *Triangle) error {
	tbuf := make([]byte, 50)
	n := 0
	for n < 50 {
		l, readErr := r.Read(tbuf[n:])
		n += l
		if readErr != nil {
			if readErr == io.EOF {
				return ErrUnexpectedEOF
			}
			return readErr
		}
	}

	offset := 0
	readBinaryPoint(tbuf, &offset, &(t.Normal))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[0]))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[1]))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[2]))
	t.Attributes = readBinaryUint16(tbuf, &offset)
	return nil
}

func readBinaryPoint(buf []byte, offset *int, p *Vec3) {
	p[0] = readBinaryFloat32(buf, offset)
	p[1] = readBinaryFloat32(buf, offset)
	p[2] = readBinaryFloat32(buf, offset)
}

func readBinaryFloat32(buf []byte, offset *int) float32 {
	v := binary.LittleEndian.Uint32(buf[*offset : (*offset)+4])
	(*offset) += 4
	return math.Float32frombits(v)
}

func readBinaryUint16(buf []byte, offset *int) uint16 {
	v := binary.LittleEndian.Uint16(buf[*offset : (*offset)+2])
	(*offset) += 2
	return v
}
