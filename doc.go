/*
Package stl implements functions to read, write, and transform files
in the Stereolithography/Surface Tesselation Language (.stl) file format
used in 3D modelling.

The format specification was taken from http://www.ennex.com/~fabbers/StL.asp,
found at http://en.wikipedia.org/wiki/STL_%28file_format%29.

While STL stores the data in single precision 32 bit floating point numbers,
the stl package does all calculations beyond simple addition in double
precision 64 bit (float64).

Usage Example

    // Read STL file
    solid, errRead := stl.ReadFile(inputFilename)
    if errRead != nil {
      fmt.Fprintln(os.Stderr, errRead)
      os.Exit(1)
    }

    // Convert from Inches to mm
    solid.Scale(25.4)

    // Write STL file
    errWrite := solid.WriteFile(outputFilename)
    if errWrite != nil {
      fmt.Fprintln(os.Stderr, errWrite)
      os.Exit(2)
    }

Everything that operates on a model is defined as a method of Solid.

Note that The STL format has two variants, a human-readable ASCII variant,
and a more compact and precise binary variant which is preferrable.

ASCII Format Specialities

The Solid.BinaryHeader field and the Triangle.Attributes fields will
be empty, after reading, as these are not part of the ASCII format. The Solid.Name
field is read from the first line after "solid ". It is not checked
against the name at the end of the file after "endsolid ". The stl package
will also not cope with Unicode byte order marks, which some text editors
might automatically place at the beginning of a file.

Binary Format Specialities

The Solid.BinaryHeader field is filled with all 80 bytes of header data.
Then, ReadFile will try to fill solid.Name with an ASCII string
read from the header data from the first byte until a \0 or a non-ASCII
character is detected.

Numerical Errors

As always when you do linear transformations on floating point numbers,
you get numerical errors. So you should expect a vertex being
rotated for 360° not to end up at exactly the original coordinates, but instead
just very close to them. As the error is usually far smaller than the available
precision of 3D printing applications, this is not an issue in most cases.
*/
package stl
