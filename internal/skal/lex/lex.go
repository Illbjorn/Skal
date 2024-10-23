package lex

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func Lex(path, in string) *token.Collection {
	tc := token.NewCollection(path, in)
	l := newLexer(path, in, tc)
	ch := tc.InputStream()

	newToken := func() token.Token {
		token := token.NewToken(tc)

		token.SetFile(l.file)
		token.SetSrcStart(l.pos)
		token.SetColumnStart(l.col)
		token.SetLineStart(l.line)

		return token
	}

	sendToken := func(token token.Token) {
		token.SetSrcEnd(l.pos)
		token.SetColumnEnd(l.col)
		token.SetLineEnd(l.line)

		ch <- token
	}

	for {
		c := l.LA()
		switch {
		// StrL
		case c == '\'' || c == '"':
			tk := newToken()
			tk.SetType(token.StrL)
			tk = eatString(tk, l)
			sendToken(tk)

		// IntL
		case isNum(c):
			tk := newToken()
			tk.SetType(token.IntL)
			tk = eatNum(tk, l)
			sendToken(tk)

		// Keyword | ID
		case isAlpha(c):
			tk := newToken()
			tk = eatWord(tk, l)
			tk.SetType(classifyKeyword(tk.Value()))

			// Discard import lines.
			if tk.Type() == token.Import {
				eatImport(l)
				break
			}
			sendToken(tk)

		// Symbol
		case isSymbol(c):
			tk := newToken()
			tk = eatSymbol(tk, l)

			// Nil == consumed line comment.
			if tk == nil {
				break
			}
			sendToken(tk)

		// Discard whitespace.
		case c == ' ', c == '\t', c == '\n':
			l.Adv()

		// Break when we hit the EOF.
		case c == rEOF:
			// Wait for the input stream to empty its queue.
			close(ch)
			tc.Wait()
			return tc

		default:
			sklog.CFatalF(
				"Fell through lexer switch statement with rune: '{rune}'.",
				"rune", string(c),
			)
		}
	}
}

func eatWord(tk token.Token, l *lexer) token.Token {
	for {
		if l.LA() == rEOF || !isAlphaNum(l.LA()) {
			return tk
		}

		tk.AddRune(l.Adv())
	}
}

func classifyKeyword(s string) string {
	switch s {
	// Keywords
	case token.New,
		token.Pub,
		token.Let,
		token.In,
		token.For,
		token.If,
		token.Elif,
		token.Else,
		token.Ret,
		token.This,
		token.Fn,
		token.Enum,
		token.Struct,
		token.True,
		token.False,
		token.Import,
		token.Defer,
		token.Extern,
		token.As,
		token.Nil,
		token.Int,
		token.Bool,
		token.Str:
		return s
	// ID
	default:
		return token.ID
	}
}

func eatNum(tk token.Token, l *lexer) token.Token {
	// Consume the first number in the series.
	tk.AddRune(l.Adv())

	// Consume numbers until we reach the end.
	for isNum(l.LA()) || (isNum(l.Cur()) && l.LA() == '.') {
		tk.AddRune(l.Adv())
	}

	return tk
}

func eatString(tk token.Token, l *lexer) token.Token {
	// Consume the left quote.
	term := l.Adv()

	for l.LA() != term {
		c := l.Adv()
		tk.AddRune(c)

		if c == rEOF {
			sklog.CFatalF(
				"Reached EOL looking for matching '{term}' in string literal.",
				"term", string(term),
			)
		}
	}

	// Consume the right quote.
	l.Adv()

	return tk
}

func eatSymbol(tk token.Token, l *lexer) token.Token {
	for {
		// Consume the next rune.
		c := l.Adv()
		tk.AddRune(c)

		// Attempt to classify it.
		switch tk.Value() {
		// --------------------------------------------------------------------------
		// Multi-byte symbols.
		// '->'
		case token.Arrow:
			tk.SetType(token.Arrow)

		// '...'
		case token.Spread:
			tk.SetType(token.Spread)

		// '..'
		case token.Concat:
			// '...'
			if l.LA() == '.' {
				break
			}

			// '..'
			tk.SetType(token.Concat)

		// '[]'
		case token.List:
			tk.SetType(tk.Value())

		// '>=' | '<=' | '!=' | '=='
		case token.GE, token.LE, token.NE, token.EQEQ:
			tk.SetType(tk.Value())

		// '&&' | '||'
		case token.And, token.Or:
			tk.SetType(tk.Value())

		// '#'
		case token.Comment:
			eatComment(l)
			return nil

		// --------------------------------------------------------------------------
		// Exclusively single-byte symbols.
		// '(' | ')'
		case token.ParenOpen, token.ParenClose:
			tk.SetType(tk.Value())

		// '{'
		case token.BraceOpen:
			tk.SetType(tk.Value())

		// '}'
		case token.BraceClose:
			tk.SetType(token.BraceClose)

		// '+'
		case token.Plus:
			tk.SetType(token.Plus)

		// '*'
		case token.Mult:
			tk.SetType(token.Mult)

		// ','
		case token.Comma:
			tk.SetType(token.Comma)

		// ';'
		case token.SemiColon:
			tk.SetType(token.SemiColon)

		// ']'
		// The open bracket is further down since it can indicate an empty list
		// literal ('[]').
		case token.BrackClose:
			tk.SetType(tk.Value())

		// ':'
		case token.Colon:
			tk.SetType(token.Colon)

		// --------------------------------------------------------------------------
		// Multi-byte symbol passthrough.
		// This just prevents exclusively multi-byte symbols from hitting the default
		// switch case.
		// '&&'
		case "&":

		// '||'
		case "|":

		// --------------------------------------------------------------------------
		// Maybe multi-byte symbols.
		// '-' | '->'
		case token.Minus:
			// '->'
			if l.LA() == '>' {
				break
			}
			tk.SetType(token.Minus)

		// '/' | '//'
		case token.Div:
			tk.SetType(token.Div)

		// '>' | '>='
		case token.GT:
			// '>='
			if l.LA() == '=' {
				break
			}
			tk.SetType(token.GT)

		// '<' | '<='
		case token.LT:
			// '<='
			if l.LA() == '=' {
				break
			}
			tk.SetType(token.LT)

		// '!' | '!='
		case token.Not:
			// '!='
			if l.LA() == '=' {
				break
			}
			// '!'
			tk.SetType(token.Not)

		// '.' | '..' | '...'
		case token.Dot:
			// '..' | '...'
			if l.LA() == '.' {
				break
			}
			// '.'
			tk.SetType(token.Dot)

		// '==' | '='
		case token.EQ:
			// '=='
			if l.LA() == '=' {
				break
			}

			// '='
			tk.SetType(token.EQ)

		// '[' | '[]'
		case token.BrackOpen:
			// '[]'
			if l.LA() == ']' {
				break
			}

			// '['
			tk.SetType(token.BrackOpen)
		default:
			sklog.CFatalF(
				"Found unknown symbol: '{symbol}'.",
				"symbol", tk.Value(),
			)
		}

		if tk.Type() != "" {
			return tk
		}
	}
}

func eatImport(l *lexer) {
	for l.LA() != '\n' && l.LA() != rEOF {
		l.Adv()
	}
}

func eatComment(l *lexer) {
	for {
		if l.LA() == '\n' {
			return
		}

		l.Adv()
	}
}
