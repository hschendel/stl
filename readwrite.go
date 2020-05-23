package stl

// This file defines the top level reading and writing operations
// for the stl package

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

// ErrIncompleteBinaryHeader is used when reading binary STL files with incomplete header.
var ErrIncompleteBinaryHeader = errors.New("Incomplete STL binary header, 84 bytes expected")

// ErrUnexpectedEOF is used by ReadFile and ReadAll to signify an incomplete file.
var ErrUnexpectedEOF = errors.New("Unexpected end of file")

var asciiStart = []byte("solid ")

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
// STL binary format, beginning with a 84 byte header.
// It also supports binary files which also start with "solid ".
// This should not be the case, but there are such files out there.
func ReadAll(file *os.File) (solid *Solid, err error) {
	isASCII, isASCIIErr := isASCIIFile(file.Name())
	if isASCIIErr != nil {
		err = isASCIIErr
		return
	}

	if isASCII {
		solid, err = readAllASCII(bufio.NewReader(file))
		if solid != nil {
			solid.IsAscii = true
		}
	} else {
		solid, err = readAllBinary(bufio.NewReader(file))
		// if read was successful, solid.IsAscii will be initialized to false
	}

	return
}

// isASCIIFile detects if the file is in STL ASCII format or if it is binary otherwise.
func isASCIIFile(fileName string) (isASCII bool, err error) {
	file, openErr := os.Open(fileName)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return
	}

	r := bufio.NewReader(file)

	const headerSize = 84
	first84 := make([]byte, headerSize)
	isASCII = false
	n, readErr := r.Read(first84)
	if n < 6 || readErr == io.EOF {
		err = ErrUnexpectedEOF
		return
	} else if readErr != nil {
		err = readErr
		return
	}

	first6 := first84[:6] // "solid "

	if bytes.Equal(first6, asciiStart) {
		isASCII = true

		// Some binary files out there also start with "solid ".
		// To check these files, read the length field at offset 80.
		// It specifies the number of triangles in the file.
		// If it matches the actual length of the file, it is most likely binary.
		// https://stackoverflow.com/questions/7377954/how-to-detect-that-this-is-a-valid-valid-binary-stlstereolithography-file/7394842#7394842
		length := binary.LittleEndian.Uint32(first84[80:84])
		if int64(length)*50+headerSize == fileInfo.Size() {
			isASCII = false
		}
	}

	return
}

// WriteFile creates file with name filename and write contents of this Solid.
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
	}
	return bufWriter.Flush()
}

// WriteAll writes the contents of this solid to an io.Writer. Depending on solid.IsAscii
// the STL ASCII format, or the STL binary format is used. If IsAscii
// is false, and the binary format is used, solid.Name will be used for
// the header, if solid.BinaryHeader is empty.
func (solid *Solid) WriteAll(w io.Writer) error {
	if solid.IsAscii {
		return writeSolidASCII(w, solid)
	}
	return writeSolidBinary(w, solid)
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
