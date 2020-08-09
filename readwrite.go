package stl

// This file defines the top level reading and writing operations
// for the stl package

import (
	"bufio"
	"errors"
	"io"
	"os"
)

// ErrIncompleteBinaryHeader is used when reading binary STL files with incomplete header.
var ErrIncompleteBinaryHeader = errors.New("incomplete STL binary header, 84 bytes expected")

// ErrUnexpectedEOF is used by ReadFile and ReadAll to signify an incomplete file.
var ErrUnexpectedEOF = errors.New("unexpected end of file")

// ReadFile reads the contents of a file into a new Solid object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Shorthand for os.Open and ReadAll
func ReadFile(filename string) (solid *Solid, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	return ReadAll(file)
}

// ReadAll reads the contents of a file into a new Solid object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Because of this,
// the file pointer has to be at the beginning of the file.
func ReadAll(r io.ReadSeeker) (solid *Solid, err error) {
	isBinary, err := isBinaryFile(r)
	if err != nil {
		return
	}
	if _, err = r.Seek(0, io.SeekStart); err != nil {
		return
	}

	if isBinary {
		solid, err = readAllBinary(r)
		// if read was successful, solid.IsAscii will be initialized to false
	} else {
		solid, err = readAllASCII(r)
		if solid != nil {
			solid.IsAscii = true
		}
	}

	return
}

// isBinaryFile returns true if the seekable stream tests as a binary file by
// matching triangle count (in header) and file size
func isBinaryFile(r io.ReadSeeker) (isBinary bool, err error) {
	var header [binaryHeaderSize]byte
	_, err = r.Read(header[:])
	if err != nil {
		if err == io.EOF { // too short to meet spec
			err = nil
		}
		return
	}
	triangleCount := triangleCountFromBinaryHeader(header[:])
	expectedFileLength := int64(triangleCount)*binaryTriangleSize + binaryHeaderSize
	actualFileLength, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return
	}
	isBinary = expectedFileLength == actualFileLength
	return
}

// WriteFile creates file with name filename and write contents of this Solid.
// Shorthand for os.Create and Solid.WriteAll
func (s *Solid) WriteFile(filename string) error {
	file, createErr := os.Create(filename)
	if createErr != nil {
		return createErr
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	err := s.WriteAll(bufWriter)
	if err != nil {
		return err
	}
	return bufWriter.Flush()
}

// WriteAll writes the contents of this solid to an io.Writer. Depending on solid.IsAscii
// the STL ASCII format, or the STL binary format is used. If IsAscii
// is false, and the binary format is used, solid.Name will be used for
// the header, if solid.BinaryHeader is empty.
func (s *Solid) WriteAll(w io.Writer) error {
	if s.IsAscii {
		return writeSolidASCII(w, s)
	}
	return writeSolidBinary(w, s)
}

// Extracts an ASCII string from a byte slice. Reads all characters
// from the beginning until a \0 or a non-ASCII character is found.
func extractASCIIString(byteData []byte) string {
	i := 0
	for i < len(byteData) && byteData[i] < byte(128) && byteData[i] != byte(0) {
		i++
	}
	return string(byteData[0:i])
}
