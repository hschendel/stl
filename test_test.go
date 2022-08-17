package stl

// This file contains code used by multiple test cases, and tests
// for the test code.

import (
	"bytes"
	"io/ioutil"
	"testing"
)

func makeTestSolid() *Solid {
	return &Solid{
		Name:    "Simple",
		IsAscii: true,
		Triangles: []Triangle{
			{
				Normal: Vec3{0, 0, -1},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{0, 1, 0},
					{1, 0, 0},
				},
			},
			{
				Normal: Vec3{0, -1, 0},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{1, 0, 0},
					{0, 0, 1},
				},
			},
			{
				Normal: Vec3{0.57735, 0.57735, 0.57735},
				Vertices: [3]Vec3{
					{0, 0, 1},
					{1, 0, 0},
					{0, 1, 0},
				},
			},
			{
				Normal: Vec3{-1, 0, 0},
				Vertices: [3]Vec3{
					{0, 0, 0},
					{0, 0, 1},
					{0, 1, 0},
				},
			},
		},
	}
}

func TestSolidSameOrderEqual(t *testing.T) {
	testSolid := makeTestSolid()
	if !testSolid.sameOrderAlmostEqual(testSolid) {
		t.Fatal("self comparison failed")
	}
}

func cmpFiles(filename1, filename2 string) (eq bool, err error) {
	data1, err1 := ioutil.ReadFile(filename1)
	if err1 != nil {
		err = err1
		return
	}
	data2, err2 := ioutil.ReadFile(filename2)
	if err2 != nil {
		err = err2
		return
	}
	eq = bytes.Equal(data1, data2)
	return
}

// true if s and o are identical up to the 5th decimal
func (s *Solid) sameOrderAlmostEqual(o *Solid) bool {
	if !(bytes.Equal(s.BinaryHeader, o.BinaryHeader) &&
		s.Name == o.Name &&
		len(s.Triangles) == len(o.Triangles) &&
		s.IsAscii == o.IsAscii) {
		return false
	}
	for i, t := range s.Triangles {
		if !t.sameOrderAlmostEqual(&o.Triangles[i], 0.000001) {
			return false
		}
	}
	return true
}

func (t *Triangle) sameOrderAlmostEqual(o *Triangle, tol float32) bool {
	return t.Normal.AlmostEqual(o.Normal, tol) &&
		t.Vertices[0].AlmostEqual(o.Vertices[0], tol) &&
		t.Vertices[1].AlmostEqual(o.Vertices[1], tol) &&
		t.Vertices[2].AlmostEqual(o.Vertices[2], tol) &&
		t.Attributes == o.Attributes
}
