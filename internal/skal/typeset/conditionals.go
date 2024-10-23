package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

/*------------------------------------------------------------------------------
 * If
 *----------------------------------------------------------------------------*/

func NewIf(n *parse.Node, p SkalType) If {
	return If{SkalType: NewBase(n, p)}
}

type If struct {
	SkalType
	Else       *Else
	Conditions []*Value
	Block      []*Statement
	Elifs      []*Elif
}

func buildIf(n node, p SkalType) If {
	s := NewIf(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Conditions
		case token.Conditions:
			s.Conditions = buildConditions(s.Conditions, child, &s)

		// Elifs
		case token.Elif:
			elif := buildElif(child, &s)
			s.Elifs = append(s.Elifs, &elif)

		// Else
		case token.Else:
			nelse := buildElse(child, &s)
			s.Else = &nelse

		// Block
		case token.Block:
			s.Block = append(s.Block, buildBlock(child, &s)...)

		default:
			sklog.UnexpectedType("typeset if node", child.Type)
		}
	}

	return s
}

/*------------------------------------------------------------------------------
 * Elif
 *----------------------------------------------------------------------------*/

func NewElif(n *parse.Node, p SkalType) Elif {
	return Elif{SkalType: NewBase(n, p)}
}

type Elif struct {
	SkalType
	Conditions []*Value
	Block      []*Statement
}

func buildElif(n node, p SkalType) Elif {
	e := NewElif(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Conditions
		case token.Conditions:
			e.Conditions = buildConditions(e.Conditions, child, &e)

		// Block
		case token.Block:
			e.Block = append(e.Block, buildBlock(child, &e)...)

		default:
			sklog.UnexpectedType("typeset elif node", child.Type)
		}
	}

	return e
}

/*------------------------------------------------------------------------------
 * Else
 *----------------------------------------------------------------------------*/

func NewElse(n *parse.Node, p SkalType) Else {
	return Else{SkalType: NewBase(n, p)}
}

type Else struct {
	SkalType
	Block []*Statement
}

func buildElse(n node, p SkalType) Else {
	e := NewElse(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Block
		case token.Block:
			e.Block = append(e.Block, buildBlock(child, &e)...)

		default:
			sklog.UnexpectedType("typeset else node", child.Type)
		}
	}

	return e
}

/*------------------------------------------------------------------------------
 * Conditions
 *----------------------------------------------------------------------------*/

func buildConditions(vs []*Value, n node, p SkalType) []*Value {
	for _, child := range n.Children {
		switch child.Type {
		// Values
		case token.Value:
			value := buildValue(child, p)
			vs = append(vs, &value)

		default:
			sklog.UnexpectedType("typeset conditions node", child.Type)
		}
	}

	return vs
}
