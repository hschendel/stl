stl
===

A library to read, write, and transform Stereolithography (.stl) files in Go.
It is used in the command line STL manipulation tool [stltool](https://github.com/hschendel/stltool).

Features
--------

* Read and write STL files in either binary or ASCII form
* Check correctness of STL files
* Measure models
* Various linear model transformations
  * Scale
  * Rotate
  * Translate (Move)
  * Fit into box
  * Apply generic 4x4 transformation matrix

Applications
------------

* Save 3D models as STL
* Import STL models
* Repair and manipulation of STL models
* General pre-processing before processing STL models in a 3D printing slicer
* Writing a slicer in Go (I hope someone does this one day)

Installation
------------

Using go's builtin installation mechanism:

    go get github.com/hschendel/stl

Usage Example
-------------

```go
solid, errRead := stl.ReadFile(inputFilename)
if errRead != nil {
  // handle
}
solid.Scale(25.4) // Convert from Inches to mm
errWrite := solid.WriteFile(outputFilename)
```

Further Reading
---------------

The [package godoc documentation](https://godoc.org/github.com/hschendel/stl)
should be helpful. Or just start a local godoc
server using this command:

    godoc -http=:6060

Then open http://localhost:6060/pkg/github.com/hschendel/stl/ in your browser.

License
-------

The stl package is licensed under the MIT license. See LICENSE file.
