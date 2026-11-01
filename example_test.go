package stl_test

import (
	"fmt"
	"os"

	"github.com/hschendel/stl"
)

// ExampleSolid_Stretch demonstrates how to read, transform, and write STL.
func ExampleSolid_Stretch() {
	solid, err := stl.ReadFile("testdata/simple_bin.stl")
	if err != nil {
		panic(err)
	}

	solid.Stretch(stl.Vec3{1, 2, 1})

	solid.IsAscii = true
	solid.Name = "Transformed"

	if err := solid.WriteAll(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// solid Transformed
	// facet normal 0 0 -1
	//   outer loop
	//     vertex 0 0 0
	//     vertex 0 2 0
	//     vertex 1 0 0
	//   endloop
	// endfacet
	// facet normal 0 -1 0
	//   outer loop
	//     vertex 0 0 0
	//     vertex 1 0 0
	//     vertex 0 0 1
	//   endloop
	// endfacet
	// facet normal 0.6666667 0.33333334 0.6666667
	//   outer loop
	//     vertex 0 0 1
	//     vertex 1 0 0
	//     vertex 0 2 0
	//   endloop
	// endfacet
	// facet normal -1 0 0
	//   outer loop
	//     vertex 0 0 0
	//     vertex 0 0 1
	//     vertex 0 2 0
	//   endloop
	// endfacet
	// endsolid Transformed
}

// statsCollector is an example for a custom data structure that implements stl.Writer for "stream processing"
// of STL data.
type statsCollector struct {
	TriangleCount uint32
	TotalSurface  float64
	Name          string
}

func (c *statsCollector) SetName(name string) {
	c.Name = name
}

func (c *statsCollector) SetBinaryHeader(header []byte) {}

func (c *statsCollector) SetASCII(bool) {}

func (c *statsCollector) SetTriangleCount(n uint32) {
	// We do not use this, as the example loads an ASCII file where this is not called.
}

func (c *statsCollector) AppendTriangle(t stl.Triangle) {
	c.TriangleCount++

	u := t.Vertices[1].Diff(t.Vertices[0])
	v := t.Vertices[2].Diff(t.Vertices[0])
	surface := u.Cross(v).Len() * 0.5
	c.TotalSurface += surface
}

// ExampleCopyFile demonstrates how to use the stl.Writer interface in order to build a custom data structure
// while processing STL input.
func ExampleCopyFile() {
	var stats statsCollector
	if err := stl.CopyFile("testdata/simple_ascii.stl", &stats); err != nil {
		panic(err)
	}

	fmt.Println(stats.Name)
	fmt.Println(stats.TriangleCount)
	fmt.Println(stats.TotalSurface)

	// Output:
	// Simple
	// 4
	// 2.3660254037844384
}
