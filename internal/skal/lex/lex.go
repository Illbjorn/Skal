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

func classifyKeyword(s string) token.Type {
	switch s {
	case token.New.String():
		return token.New
	case token.New.String():
		return token.New
	case token.Pub.String():
		return token.Pub
	case token.Let.String():
		return token.Let
	case token.In.String():
		return token.In
	case token.For.String():
		return token.For
	case token.If.String():
		return token.If
	case token.Elif.String():
		return token.Elif
	case token.Else.String():
		return token.Else
	case token.Ret.String():
		return token.Ret
	case token.This.String():
		return token.This
	case token.Fn.String():
		return token.Fn
	case token.Enum.String():
		return token.Enum
	case token.Struct.String():
		return token.Struct
	case token.True.String():
		return token.True
	case token.False.String():
		return token.False
	case token.Import.String():
		return token.Import
	case token.Defer.String():
		return token.Defer
	case token.Extern.String():
		return token.Extern
	case token.As.String():
		return token.As
	case token.Nil.String():
		return token.Nil
	case token.Int.String():
		return token.Int
	case token.Bool.String():
		return token.Bool
	case token.Str.String():
		return token.Str
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
		case token.Arrow.String():
			tk.SetType(token.Arrow)

		// '...'
		case token.Spread.String():
			tk.SetType(token.Spread)

		// '..'
		case token.Concat.String():
			// '...'
			if l.LA() == '.' {
				break
			}

			// '..'
			tk.SetType(token.Concat)

		// '[]'
		case token.List.String():
			tk.SetType(token.List)

		// '>='
		case token.GE.String():
			tk.SetType(token.GE)

		// '<='
		case token.LE.String():
			tk.SetType(token.LE)

		// '!='
		case token.NE.String():
			tk.SetType(token.NE)

		// '=='
		case token.EQEQ.String():
			tk.SetType(token.EQEQ)

		// '&&'
		case token.And.String():
			tk.SetType(token.And)

		// '||'
		case token.Or.String():
			tk.SetType(token.Or)

		// '#'
		case token.Comment.String():
			eatComment(l)
			return nil

		// --------------------------------------------------------------------------
		// Exclusively single-byte symbols.
		// '(' | ')'
		case token.ParenOpen.String():
			tk.SetType(token.ParenOpen)

		case token.ParenClose.String():
			tk.SetType(token.ParenClose)

		// '{'
		case token.BraceOpen.String():
			tk.SetType(token.BraceOpen)

		// '}'
		case token.BraceClose.String():
			tk.SetType(token.BraceClose)

		// '+'
		case token.Plus.String():
			tk.SetType(token.Plus)

		// '*'
		case token.Mult.String():
			tk.SetType(token.Mult)

		// ','
		case token.Comma.String():
			tk.SetType(token.Comma)

		// ';'
		case token.SemiColon.String():
			tk.SetType(token.SemiColon)

		// ']'
		// The open bracket is further down since it can indicate an empty list
		// literal ('[]').
		case token.BrackClose.String():
			tk.SetType(token.BrackClose)

		// ':'
		case token.Colon.String():
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
		case token.Minus.String():
			// '->'
			if l.LA() == '>' {
				break
			}
			tk.SetType(token.Minus)

		// '/' | '//'
		case token.Div.String():
			tk.SetType(token.Div)

		// '>' | '>='
		case token.GT.String():
			// '>='
			if l.LA() == '=' {
				break
			}
			tk.SetType(token.GT)

		// '<' | '<='
		case token.LT.String():
			// '<='
			if l.LA() == '=' {
				break
			}
			tk.SetType(token.LT)

		// '!' | '!='
		case token.Not.String():
			// '!='
			if l.LA() == '=' {
				break
			}
			// '!'
			tk.SetType(token.Not)

		// '.' | '..' | '...'
		case token.Dot.String():
			// '..' | '...'
			if l.LA() == '.' {
				break
			}
			// '.'
			tk.SetType(token.Dot)

		// '==' | '='
		case token.EQ.String():
			// '=='
			if l.LA() == '=' {
				break
			}

			// '='
			tk.SetType(token.EQ)

		// '[' | '[]'
		case token.BrackOpen.String():
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

		if tk.Type() > 0 {
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
