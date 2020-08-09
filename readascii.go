package stl

// This file defines a parser for the STL ASCII format.

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
)

func readAllASCII(r io.Reader) (solid *Solid, err error) {
	var sd Solid
	p := newParser(r)
	if p.Parse(&sd) {
		solid = &sd
	} else {
		err = errors.New(p.ErrorText)
	}
	return
}

type parser struct {
	line             int
	errors           *list.List
	currentWord      string
	currentLine      []byte
	eof              bool
	lineScanner      *bufio.Scanner
	wordScanner      *bufio.Scanner
	HeaderError      bool
	TrianglesSkipped bool
	ErrorText        string
}

func newParser(reader io.Reader) *parser {
	var p parser
	p.errors = list.New()
	p.eof = false
	p.lineScanner = bufio.NewScanner(reader)
	p.nextLine()
	return &p
}

func (p *parser) addError(msg string) {
	p.errors.PushBack(fmt.Sprintf("%d: %s", p.line, msg))
}

const (
	idNone  = 0
	idSolid = 1 << iota
	idFacet
	idNormal
	idOuter
	idLoop
	idVertex
	idEndloop
	idEndfacet
	idEndsolid
)

var identRegexps = map[int]*regexp.Regexp{
	idSolid:                regexp.MustCompile("^solid$"),
	idFacet:                regexp.MustCompile("^facet$"),
	idNormal:               regexp.MustCompile("^normal$"),
	idOuter:                regexp.MustCompile("^outer$"),
	idLoop:                 regexp.MustCompile("^loop$"),
	idVertex:               regexp.MustCompile("^vertex$"),
	idEndloop:              regexp.MustCompile("^endloop$"),
	idEndfacet:             regexp.MustCompile("^endfacet$"),
	idEndsolid:             regexp.MustCompile("^endsolid$"),
	(idFacet | idEndsolid): regexp.MustCompile(`^(facet|endsolid)$`),
}

var reFloat = regexp.MustCompile(`^[+-]?\d+(\.\d+)?([eE][+-]?\d+)?$`)

var idents = map[int]string{
	idSolid:    "solid",
	idFacet:    "facet",
	idNormal:   "normal",
	idOuter:    "outer",
	idLoop:     "loop",
	idVertex:   "vertex",
	idEndloop:  "endloop",
	idEndfacet: "endfacet",
	idEndsolid: "endsolid",
}

func (p *parser) Parse(solid *Solid) bool {
	if p.eof {
		p.HeaderError = true
		p.addError("File is empty")
	} else {
		p.HeaderError = !p.parseASCIIHeaderLine(solid)

		triangles := list.New()
	TriangleLoop:
		for !p.eof && !p.isCurrentTokenIdent(idEndsolid) {
			if !p.isCurrentTokenIdent(idFacet) {
				p.addError(`"facet" or "endsolid" expected`)
				switch p.skipToToken(idFacet | idEndsolid) {
				case idEndsolid, idNone:
					break TriangleLoop
				}
			}

			var t Triangle
			if p.parseFacet(&t) {
				triangles.PushBack(&t)
			} else {
				p.TrianglesSkipped = true
				p.skipToToken(idFacet | idEndsolid)
			}
		}
		solid.Triangles = make([]Triangle, triangles.Len())
		for i, e := 0, triangles.Front(); e != nil; e = e.Next() {
			solid.Triangles[i] = *((e.Value).(*Triangle))
			i++
		}
	}

	success := !p.HeaderError && !p.TrianglesSkipped && p.consumeToken(idEndsolid)
	p.generateErrorText()
	return success
}

func (p *parser) generateErrorText() {
	var buf bytes.Buffer
	if p.TrianglesSkipped {
		buf.WriteString("Triangles had to be skipped.\n")
	}
	for e := p.errors.Front(); e != nil; e = e.Next() {
		buf.WriteString(e.Value.(string))
		buf.WriteString("\n")
	}
	p.ErrorText = buf.String()
}

var expectedASCIIHeaderPrefix = []byte("solid ")

func (p *parser) parseASCIIHeaderLine(solid *Solid) bool {
	var success bool
	if p.eof {
		p.addError("unexpected end of file")
		success = false
	} else {
		if !bytes.HasPrefix(p.currentLine, expectedASCIIHeaderPrefix) {
			p.addError("ASCII header must start with \"solid \"")
			success = false
		} else {
			solid.Name = extractASCIIString(p.currentLine[len(expectedASCIIHeaderPrefix):])
			success = true
		}
	}
	p.nextLine()
	return success
}

func (p *parser) parseFacet(t *Triangle) bool {
	return p.consumeToken(idFacet) &&
		p.consumeToken(idNormal) && p.parsePoint(&(t.Normal)) &&
		p.consumeToken(idOuter) && p.consumeToken(idLoop) &&
		p.consumeToken(idVertex) && p.parsePoint(&(t.Vertices[0])) &&
		p.consumeToken(idVertex) && p.parsePoint(&(t.Vertices[1])) &&
		p.consumeToken(idVertex) && p.parsePoint(&(t.Vertices[2])) &&
		p.consumeToken(idEndloop) &&
		p.consumeToken(idEndfacet)
}

func (p *parser) parsePoint(pt *Vec3) bool {
	return p.parseFloat32(&(pt[0])) &&
		p.parseFloat32(&(pt[1])) &&
		p.parseFloat32(&(pt[2]))
}

func (p *parser) parseFloat32(f *float32) bool {
	if p.eof {
		return false
	}
	f64, err := strconv.ParseFloat(p.currentWord, 32)
	if err != nil {
		p.addError("Unable to parse float")
		return false
	}

	*f = float32(f64)
	p.nextWord()
	return true
}

func (p *parser) isCurrentTokenIdent(ident int) bool {
	re := identRegexps[ident]
	return re.MatchString(p.currentWord)
}

func (p *parser) skipToToken(ident int) int {
	re := identRegexps[ident]
	for { // terminates when no more next words are there, or ident has been found
		if re.MatchString(p.currentWord) {
			if ident == (idFacet | idEndsolid) {
				if identRegexps[idFacet].MatchString(p.currentWord) {
					return idFacet
				}
				return idEndsolid
			}
			return ident
		}
		if !p.nextWord() {
			return idNone
		}
	}
}

func (p *parser) consumeToken(ident int) bool {
	re := identRegexps[ident]
	if !re.MatchString(p.currentWord) {
		ident := idents[ident]
		p.addError("\"" + ident + "\" expected")
		return false
	}

	p.nextWord()
	return true
}

func (p *parser) nextWord() bool {
	if p.eof {
		return false
	}
	// Try to advance word scanner
	if p.wordScanner.Scan() {
		p.currentWord = p.wordScanner.Text()
		return true
	}
	if p.wordScanner.Err() == nil { // line has ended
		return p.nextLine()
	}
	p.addError(p.wordScanner.Err().Error())
	p.currentLine = nil
	p.currentWord = ""
	p.eof = true
	return false
}

func (p *parser) nextLine() bool {
	if p.lineScanner.Scan() {
		p.currentLine = p.lineScanner.Bytes()
		p.line++
		p.wordScanner = bufio.NewScanner(bytes.NewReader(p.currentLine))
		p.wordScanner.Split(bufio.ScanWords)
		return p.nextWord()
	}

	if p.lineScanner.Err() != nil {
		p.addError(p.lineScanner.Err().Error())
	}
	p.currentLine = nil
	p.currentWord = ""
	p.eof = true
	return false
}
