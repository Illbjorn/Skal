package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewFn(n *parse.Node, p SkalType) Fn {
	return Fn{SkalType: NewBase(n, p)}
}

type Fn struct {
	SkalType
	Args        []*FnArg
	Block       []*Statement
	Values      []*Value
	Method      bool
	Constructor bool
}

func NewFnArg(n *parse.Node, p SkalType) FnArg {
	return FnArg{SkalType: NewBase(n, p)}
}

type FnArg struct {
	SkalType
	Vararg bool
}

func buildFn(n node, p SkalType) Fn {
	fn := NewFn(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			fn.AddRef(child.Value)

		// Args
		case token.FnArg:
			arg := buildFnArg(child, &fn)
			fn.Args = append(fn.Args, &arg)

		// Block
		case token.Block:
			fn.Block = append(fn.Block, buildBlock(child, &fn)...)

		// Anonymous fns have values rather than statements.
		case token.Value:
			value := buildValue(child, p)
			fn.Values = append(fn.Values, &value)

		// 'new'
		// This is a special struct method name which overrides the auto-generated
		// constructor.
		case token.New:
			if _, ok := p.(*Struct); ok {
				// Indicate this fn is a constructor.
				fn.Constructor = true

				// Indicate that we should NOT generate a default constructor.
				p.(*Struct).NoConstructor = true
			}

		// Type hint.
		case token.TypeHint:
			fn.SetType(child.Value)

		default:
			sklog.UnexpectedType("typeset fn node", child.Type)
		}
	}

	return fn
}

func buildFnArg(n node, p SkalType) FnArg {
	arg := NewFnArg(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			arg.AddRef(child.Value)

		// '...'
		case token.Spread:
			arg.Vararg = true

		// Type Hint
		case token.TypeHint:
			arg.SetType(child.Value)

		default:
			sklog.UnexpectedType("typeset fn arg node", child.Type)
		}
	}

	return arg
}
