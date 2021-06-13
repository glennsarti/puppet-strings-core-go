package puppet_strings_go

import (
	"bytes"
	"unicode/utf8"
)

type StringReader interface {
	Next() (c rune, start int)
	Peek() (c rune, start int)
	NextLine() (line string, start int, end int)
	PeekNextLine() (line string, start int, end int)
	Advance(size int)
	Pos() int
	SetPos(pos int)
	SubString(start int, end int) string
	PeekUntilEnd() string
	IsEOF() bool

	// // Returns the the current rune and its size in the parsed string. The position does not change
	// Peek() (c rune, size int)

	// Advance(size int)

	// Pos() int

	// SetPos(int)

	// // Returns the string that is backing the reader
	// Text() string

	// // Returns the substring starting at start and up to, but not including, the current position
	// From(start int) string
}

type stringReader struct {
	pos int
	text string
}

func NewStringReader(text string) StringReader {
	return &stringReader{pos: 0, text: text}
}

func (sr *stringReader) Next() (c rune, start int) {
	start = sr.pos
	if sr.pos >= len(sr.text) {
		return
	}
	c = rune(sr.text[sr.pos])
	if c < utf8.RuneSelf {
		sr.pos++
		return
	}
	c, size := utf8.DecodeRuneInString(sr.text[sr.pos:])
	if c == utf8.RuneError {
		panic("invalid unicode character")
	}
	sr.pos += size
	return
}

func (sr *stringReader) Peek() (line rune, start int) {
	start = sr.pos
	c, _ := sr.Next()
	sr.pos = start
	return c, start
}

func (sr *stringReader) NextLine() (line string, start int, end int) {
	buf := bytes.NewBufferString(``)
	start = sr.pos
	for {
		ct, _ := sr.Next()
		if (ct == 0) || (ct == '\n') { return buf.String(), start, sr.pos }
		if (ct == '\r') { continue } // Ignore CR
		buf.WriteRune(ct)
	}
}

func (sr *stringReader) PeekNextLine() (line string, start int, end int) {
	start = sr.pos
	line, _, end = sr.NextLine()
	sr.pos = start
	return line, start, end
}

func (sr *stringReader) Advance(size int) {
	sr.pos = sr.pos + size
}

func (sr *stringReader) Pos() (int) {
	return sr.pos
}

func (sr *stringReader) SetPos(pos int) {
	sr.pos = pos
}

func (sr *stringReader) SubString(start int, end int) string {
	return sr.text[start:end]
}

func (sr *stringReader) PeekUntilEnd() string {
	return sr.text[sr.pos:]
}
func (sr *stringReader) IsEOF() bool {
	return sr.pos >= len(sr.text)
}
