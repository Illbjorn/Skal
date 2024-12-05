package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

/*------------------------------------------------------------------------------
 * For
 *----------------------------------------------------------------------------*/

func NewFor(n *parse.Node, p SkalType) For {
	return For{SkalType: NewBase(n, p)}
}

type For struct {
	SkalType
	ForType   string
	Iterators []*ForI
	Iterables []*ForV
	Block     []*Statement
}

func buildFor(n node, p SkalType) For {
	f := NewFor(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// 'in' | '='
		case token.ForType:
			f.ForType = child.Value

		// Iterators
		// k, v
		// i
		case token.ForIterator:
			iterator := buildForIterator(child, &f)
			f.Iterators = append(f.Iterators, &iterator)

		// Iterables
		// 1, 2, 1
		// 1, 2
		// collection
		case token.ForIterable:
			iterable := buildForIterable(child, &f)
			f.Iterables = append(f.Iterables, &iterable)

		// Block
		case token.Block:
			f.Block = append(f.Block, buildBlock(child, &f)...)

		default:
			sklog.UnexpectedType("typeset for node", child.Type.String())
		}
	}

	return f
}

/*------------------------------------------------------------------------------
 * For Iterator
 *----------------------------------------------------------------------------*/

func NewForI(n *parse.Node, p SkalType) ForI {
	return ForI{SkalType: NewBase(n, p)}
}

type ForI struct {
	SkalType
}

func buildForIterator(n node, p SkalType) ForI {
	f := NewForI(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// ID | Reference
		case token.ID, token.Ref:
			f = *buildRef(child, &f)

		default:
			sklog.UnexpectedType("typeset for iterator node", child.Type.String())
		}
	}

	return f
}

/*------------------------------------------------------------------------------
 * For Iterable
 *----------------------------------------------------------------------------*/

func NewForV(n *parse.Node, p SkalType) ForV {
	return ForV{SkalType: NewBase(n, p)}
}

// A single For iterable.
type ForV struct {
	SkalType
	Value        string
	IterableType token.Type
}

func buildForIterable(n node, p SkalType) ForV {
	f := NewForV(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Reference
		case token.Ref:
			f = *buildRef(child, &f)

		// ID
		case token.ID:
			f.IterableType = token.ID
			f.Value = child.Value

		// IntL
		case token.IntL:
			f.IterableType = token.IntL
			f.Value = child.Value

		default:
			sklog.UnexpectedType("typeset for iterable node", child.Type.String())
		}
	}

	return f
}
