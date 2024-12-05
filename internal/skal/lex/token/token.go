package token

import (
	"bytes"
	"strconv"

	"github.com/illbjorn/fstr"
)

var itoa = strconv.Itoa

type Token interface {
	SetType(t Type)
	Type() Type
	LineStart() int
	SetLineStart(i int)
	LineEnd() int
	SetLineEnd(i int)
	ColumnStart() int
	SetColumnStart(i int)
	ColumnEnd() int
	SetColumnEnd(i int)
	SetSrcStart(i int)
	SetSrcEnd(i int)
	AddRune(r rune)
	File() string
	SetFile(file string)
	Value() string
	String() string
	Src() string
	SrcLine() string
}

func NewToken(tc *Collection) Token {
	token := &token{value: bytes.NewBuffer(nil)}
	token.start = new(Position)
	token.end = new(Position)

	token.src = func() string {
		if token.end.Abs < token.start.Abs {
			return ""
		}
		return tc.src[token.start.Abs:token.end.Abs]
	}

	token.srcLine = func() string {
		return tc.SrcLine(token.start.Line)
	}

	return token
}

var _ Token = (*token)(nil)

type token struct {
	value   *bytes.Buffer
	start   *Position
	end     *Position
	src     func() string
	srcLine func() string
	file    string
	_type   Type
}

// Type
func (tk *token) SetType(t Type) { tk._type = t }
func (tk *token) Type() Type     { return tk._type }

// Line Position
func (tk *token) LineStart() int     { return tk.start.Line }
func (tk *token) SetLineStart(i int) { tk.start.Line = i }
func (tk *token) LineEnd() int       { return tk.end.Line }
func (tk *token) SetLineEnd(i int)   { tk.end.Line = i }

// Column Position
func (tk *token) ColumnStart() int     { return tk.start.Col }
func (tk *token) SetColumnStart(i int) { tk.start.Col = i }
func (tk *token) ColumnEnd() int       { return tk.end.Col }
func (tk *token) SetColumnEnd(i int)   { tk.end.Col = i }

// Absolute Source Position
func (tk *token) SetSrcStart(i int) { tk.start.Abs = i }
func (tk *token) SetSrcEnd(i int)   { tk.end.Abs = i }

// Buffer Writer
func (tk *token) AddRune(r rune) { tk.value.WriteRune(r) }

// Source File Get/Set
func (tk *token) File() string        { return tk.file }
func (tk *token) SetFile(file string) { tk.file = file }

// Source Text Getters
func (tk *token) Src() string     { return tk.src() }
func (tk *token) SrcLine() string { return tk.srcLine() }

// Buffer Stringer
func (tk *token) Value() string { return tk.value.String() }

// fmt.Stringer
func (tk *token) String() string {
	return fstr.Pairs(
		`File    : {File}
Type    : {type}
Start X : {sx}
Start Y : {sy}
End X   : {ex}
End Y   : {ey}`,
		"File", tk.file,
		"type", tk._type.String(),
		"sx", itoa(tk.start.Col),
		"sy", itoa(tk.start.Line),
		"ex", itoa(tk.end.Col),
		"ey", itoa(tk.end.Line),
	)
}

type Position struct {
	Abs, Col, Line int
}
