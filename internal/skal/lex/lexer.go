package lex

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
)

func newLexer(path, s string, tc *token.Collection) *lexer {
	return &lexer{
		file:  path,
		runes: []rune(s),
		pos:   -1,
		len:   len([]rune(s)),
		col:   1,
		line:  1,
		tc:    tc,
	}
}

type lexer struct {
	tc               *token.Collection
	file             string
	runes            []rune
	currentLineStart int
	pos              int
	line             int
	col              int
	len              int
}

// Adv moves the position index forward by one and returns the now-current rune
// in the slice.
func (l *lexer) Adv() rune {
	// Make sure we're not indexing outside the slice bounds.
	if l.pos >= l.len-1 {
		return rEOF
	}

	// Bump the pointer.
	l.pos++

	// On line feed, increment line count and reset column index.
	if l.runes[l.pos] == '\n' {
		// Add the line position.
		l.tc.MarkLine(l.line, l.currentLineStart, l.pos)

		// Increment the line count.
		l.line++

		// Reset the column position.
		l.col = 1

		// Mark the new-current-line start.
		l.currentLineStart = l.pos

		// Otherwise, just bump the column position.
	} else {
		l.col++
	}

	return l.Cur()
}

// Cur returns the current-index rune in the runes slice.
func (l *lexer) Cur() rune {
	return l.runes[l.pos]
}

// LA returns the current-index+1 rune in the runes slice.
func (l *lexer) LA() rune {
	if l.pos+1 < len(l.runes) {
		return l.runes[l.pos+1]
	}

	return rEOF
}

// LB returns the current position - 1 rune in the runes slice.
func (l *lexer) LB() rune {
	if l.pos-1 >= 0 {
		return l.runes[l.pos-1]
	}

	return rEOF
}
