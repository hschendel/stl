package stl

// Tests for reading and writing STL files.

import (
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

const testFilenameSimpleASCII = "testdata/simple_ascii.stl"
const testFilenameSimpleBinary = "testdata/simple_bin.stl"
const testFilenameComplexASCII = "testdata/complex_ascii.stl"
const testFilenameComplexBinary = "testdata/complex_bin.stl"

func TestIsAsciiFile(t *testing.T) {
	asciiFile, openASCIIErr := os.Open(testFilenameSimpleASCII)
	if openASCIIErr != nil {
		t.Fatal(openASCIIErr)
	}
	defer asciiFile.Close()

	isASCII, _, err := isASCIIFile(asciiFile)
	if err != nil {
		t.Fatal(err)
	}
	if !isASCII {
		t.Error("ASCII file not detected as ASCII")
	}

	binaryFile, openBinaryErr := os.Open(testFilenameSimpleBinary)
	if openBinaryErr != nil {
		t.Fatal(openBinaryErr)
	}
	defer binaryFile.Close()

	isASCII, _, err = isASCIIFile(binaryFile)
	if err != nil {
		t.Fatal(err)
	}
	if isASCII {
		t.Error("Binary file detected as ASCII")
	}
}

func TestReadFile_Ascii(t *testing.T) {
	solid, err := ReadFile(testFilenameSimpleASCII)
	if err != nil {
		t.Fatal("Error in ReadFile: " + err.Error())
	}
	testSolid := makeTestSolid()
	testSolid.IsAscii = true
	testSolid.BinaryHeader = nil
	if !solid.sameOrderAlmostEqual(testSolid) {
		t.Error("Not as expected")
		t.Log("Expected:\n", testSolid)
		t.Log("Found:\n", solid)
	}
}

func BenchmarkReadFile_ASCII_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ReadFile(testFilenameSimpleASCII)
		if err != nil {
			b.Fatal("Error in ReadFile: " + err.Error())
		}
	}
}

func BenchmarkReadFile_ASCII_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ReadFile(testFilenameComplexASCII)
		if err != nil {
			b.Fatal("Error in ReadFile: " + err.Error())
		}
	}
}

func TestReadFile_Binary(t *testing.T) {
	solid, err := ReadFile(testFilenameSimpleBinary)
	if err != nil {
		t.Fatal("Error in ReadFile: " + err.Error())
	}
	testSolid := makeTestSolid()
	testSolid.IsAscii = false
	testSolid.BinaryHeader = make([]byte, 80)
	copy(testSolid.BinaryHeader, string(testSolid.Name))
	if !solid.sameOrderAlmostEqual(testSolid) {
		t.Error("Not as expected")
		t.Log("Expected:\n", testSolid)
		t.Log("Found:\n", solid)
	}
}

func BenchmarkReadFile_Binary_Simple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ReadFile(testFilenameSimpleBinary)
		if err != nil {
			b.Fatal("Error in ReadFile: " + err.Error())
		}
	}
}

func BenchmarkReadFile_Binary_Complex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ReadFile(testFilenameComplexBinary)
		if err != nil {
			b.Fatal("Error in ReadFile: " + err.Error())
		}
	}
}

func TestReadAll_Binary(t *testing.T) {
	file, openErr := os.Open(testFilenameSimpleBinary)
	if openErr != nil {
		t.Fatal(openErr)
	}
	defer file.Close()

	solid, err := ReadAll(file)
	if err != nil {
		t.Fatal("Error in ReadAll: " + err.Error())
	}
	testSolid := makeTestSolid()
	testSolid.IsAscii = false
	testSolid.BinaryHeader = make([]byte, 80)
	copy(testSolid.BinaryHeader, string(testSolid.Name))
	if !solid.sameOrderAlmostEqual(testSolid) {
		t.Error("Not as expected")
		t.Log("Expected:\n", testSolid)
		t.Log("Found:\n", solid)
	}
}

func TestWriteFile_Ascii(t *testing.T) {
	tmpDirName, tmpErr := ioutil.TempDir(os.TempDir(), "stl_test")
	if tmpErr != nil {
		t.Fatal(tmpErr.Error())
	}
	defer os.RemoveAll(tmpDirName)

	tmpFileName := tmpDirName + string(os.PathSeparator) + "test_out_ascii.stl"
	testSolid := makeTestSolid()
	testSolid.IsAscii = true
	err := testSolid.WriteFile(tmpFileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	eq, cmpErr := cmpFiles(testFilenameSimpleASCII, tmpFileName)
	if cmpErr != nil {
		t.Fatal(cmpErr)
	}

	if !eq {
		t.Log("Generated file:")
		data, readErr := ioutil.ReadFile(tmpFileName)
		if readErr != nil {
			t.Fatal(readErr)
		}
		t.Log(string(data))
		t.Error("Was expected to look like " + testFilenameSimpleASCII)
	}
}

func BenchmarkWriteSmallFile_ASCII(b *testing.B) {
	tmpDirName, tmpErr := ioutil.TempDir(os.TempDir(), "stl_test")
	if tmpErr != nil {
		b.Fatal(tmpErr.Error())
	}
	defer os.RemoveAll(tmpDirName)

	testSolid := makeTestSolid()
	testSolid.IsAscii = true
	var tmpFileName string
	for i := 0; i < b.N; i++ {
		// generate a new file name for every turn, as overwriting the same file again
		// and again seems to delay some systems.
		tmpFileName = tmpDirName + string(os.PathSeparator) + strconv.Itoa(i) + "_bench_out_ascii.stl"
		err := testSolid.WriteFile(tmpFileName)
		if err != nil {
			b.Fatal("Error in WriteFile: " + err.Error())
		}
	}
}

func TestWriteFile_Binary(t *testing.T) {
	tmpDirName, tmpErr := ioutil.TempDir(os.TempDir(), "stl_test")
	if tmpErr != nil {
		t.Fatal(tmpErr.Error())
	}
	defer os.RemoveAll(tmpDirName)

	tmpFileName := tmpDirName + string(os.PathSeparator) + "test_out_binary.stl"
	testSolid := makeTestSolid()
	testSolid.IsAscii = false
	err := testSolid.WriteFile(tmpFileName)
	if err != nil {
		t.Fatal(err.Error())
	}

	eq, cmpErr := cmpFiles(testFilenameSimpleBinary, tmpFileName)
	if cmpErr != nil {
		t.Fatal(cmpErr)
	}

	if !eq {
		t.Log("Generated file:")
		data, readErr := ioutil.ReadFile(tmpFileName)
		if readErr != nil {
			t.Fatal(readErr)
		}
		t.Log(string(data))
		t.Error("Was expected to look like " + testFilenameSimpleBinary)
	}
}

func BenchmarkWriteSmallFile_Binary(b *testing.B) {
	tmpDirName, tmpErr := ioutil.TempDir(os.TempDir(), "stl_test")
	if tmpErr != nil {
		b.Fatal(tmpErr.Error())
	}
	defer os.RemoveAll(tmpDirName)

	testSolid := makeTestSolid()
	testSolid.IsAscii = false
	var tmpFileName string
	for i := 0; i < b.N; i++ {
		// generate a new file name for every turn, as overwriting the same file again
		// and again seems to delay some systems.
		tmpFileName = tmpDirName + string(os.PathSeparator) + strconv.Itoa(i) + "_bench_out_binary_small.stl"
		err := testSolid.WriteFile(tmpFileName)
		if err != nil {
			b.Fatal("Error in WriteFile: " + err.Error())
		}
	}
}

// Does assume ReadFile to work correctly
func BenchmarkWriteMediumFile_Binary(b *testing.B) {
	b.StopTimer()
	tmpDirName, tmpErr := ioutil.TempDir(os.TempDir(), "stl_test")
	if tmpErr != nil {
		b.Fatal(tmpErr)
	}
	defer os.RemoveAll(tmpDirName)

	testSolid, readErr := ReadFile(testFilenameComplexBinary)
	if readErr != nil {
		b.Fatal(readErr)
	}

	testSolid.IsAscii = false // to be sure ;-)
	b.StartTimer()
	var tmpFileName string
	for i := 0; i < b.N; i++ {
		// generate a new file name for every turn, as overwriting the same file again
		// and again seems to delay some systems.
		tmpFileName = tmpDirName + string(os.PathSeparator) + strconv.Itoa(i) + "_bench_out_binary_medium.stl"
		err := testSolid.WriteFile(tmpFileName)
		if err != nil {
			b.Fatal("Error in WriteFile: " + err.Error())
		}
	}
}
