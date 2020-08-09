package stl

// This file defines functions to emit STL ASCII files.

import (
	"fmt"
	"io"
	"strings"
)

func writeSolidASCII(w io.Writer, solid *Solid) error {
	var writeErr error
	_, writeErr = w.Write([]byte("solid " + escapeName(solid.Name)))
	if writeErr != nil {
		return writeErr
	}
	for _, triangle := range solid.Triangles {
		writeErr = writeTriangleASCII(w, &triangle)
		if writeErr != nil {
			return writeErr
		}
	}
	_, writeErr = w.Write([]byte("\nendsolid " + solid.Name + "\n"))
	if writeErr != nil {
		return writeErr
	}
	return nil
}

func escapeName(name string) string {
	name = strings.ReplaceAll(name, "\r", "\\r")
	name = strings.ReplaceAll(name, "\n", "\\n")
	return name
}

func writeTriangleASCII(w io.Writer, t *Triangle) error {
	var err error
	_, err = w.Write([]byte("\nfacet normal "))
	if err != nil {
		return err
	}
	err = writePointString(w, &(t.Normal))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n  outer loop"))
	if err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		_, err = w.Write([]byte("\n    vertex "))
		if err != nil {
			return err
		}
		err = writePointString(w, &(t.Vertices[i]))
		if err != nil {
			return err
		}
	}
	_, err = w.Write([]byte("\n  endloop\nendfacet"))
	if err != nil {
		return err
	}
	return nil
}

func writePointString(w io.Writer, p *Vec3) error {
	// %v is the easiest way I know to write floats as compact as possible
	_, err := w.Write([]byte(fmt.Sprintf("%v %v %v", p[0], p[1], p[2])))
	return err
}
