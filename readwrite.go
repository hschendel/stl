package stl

// This file defines the top level reading and writing operations
// for the stl package

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
)

// Error when reading binary STL files with incomplete header.
var ErrIncompleteBinaryHeader = errors.New("Incomplete STL binary header, 84 bytes expected")

// Used by ReadFile and ReadAll to signify an incomplete file.
var ErrUnexpectedEOF = errors.New("Unexpected end of file")

var asciiStart = []byte("solid ")

// Reads the contents of a file into a new Solid object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Shorthand for os.Open and ReadAll
func ReadFile(filename string) (solid *Solid, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	return ReadAll(bufio.NewReader(file))
}

// Reads the contents of a file into a new Solid object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Because of this,
// the file pointer has to be at the beginning of the file.
func ReadAll(r io.Reader) (solid *Solid, err error) {
	isAscii, first6, isAsciiErr := isAsciiFile(r)
	if isAsciiErr != nil {
		err = isAsciiErr
		return
	}

	if isAscii {
		solid, err = readAllAscii(r)
		if solid != nil {
			solid.IsAscii = true
		}
	} else {
		solid, err = readAllBinary(r, first6)
		// if read was successful, solid.IsAscii will be initialized to false
	}

	return
}

// Detects if the file is in STL ASCII format or if it is binary otherwise.
// It will consume 6 bytes and return them.
func isAsciiFile(r io.Reader) (isAscii bool, first6 []byte, err error) {
	first6 = make([]byte, 6) // "solid "
	isAscii = false
	n, readErr := r.Read(first6)
	if n != 6 || readErr == io.EOF {
		err = ErrUnexpectedEOF
		return
	} else if readErr != nil {
		err = readErr
		return
	}

	if bytes.Equal(first6, asciiStart) {
		isAscii = true
	}

	return
}

// Create file with name filename and write contents of this Solid.
// Shorthand for os.Create and Solid.WriteAll
func (solid *Solid) WriteFile(filename string) error {
	file, createErr := os.Create(filename)
	if createErr != nil {
		return createErr
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	err := solid.WriteAll(bufWriter)
	if err != nil {
		return err
	} else {
		return bufWriter.Flush()
	}
}

// Write contents of this solid to an io.Writer. Depending on solid.IsAscii
// the STL ASCII format, or the STL binary format is used. If IsAscii
// is false, and the binary format is used, solid.Name will be used for
// the header, if solid.BinaryHeader is empty.
func (solid *Solid) WriteAll(w io.Writer) error {
	if solid.IsAscii {
		return writeSolidAscii(w, solid)
	} else {
		return writeSolidBinary(w, solid)
	}
}

// Extracts an ASCII string from a byte slice. Reads all characters
// from the beginning until a \0 or a non-ASCII character is found.
func extractAsciiString(byteData []byte) string {
	i := 0
	for i < len(byteData) && byteData[i] < byte(128) && byteData[i] != byte(0) {
		i++
	}
	return string(byteData[0:i])
}
