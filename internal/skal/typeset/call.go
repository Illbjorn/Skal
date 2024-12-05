package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewCall(n *parse.Node, p SkalType) Call {
	return Call{SkalType: NewBase(n, p)}
}

type Call struct {
	SkalType
	Args []*CallArg
}

func NewCallArg(n *parse.Node, p SkalType) CallArg {
	return CallArg{SkalType: NewBase(n, p)}
}

type CallArg struct {
	SkalType
	Values []*Value
	Spread bool
}

func buildCall(n node, p SkalType) Call {
	c := NewCall(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Reference
		case token.Ref:
			c = *buildRef(child, &c)

		// Args
		case token.CallArg:
			arg := buildCallArg(child, &c)
			c.Args = append(c.Args, &arg)

		default:
			sklog.UnexpectedType("typeset call node", child.Type.String())
		}
	}

	return c
}

func buildCallArg(n node, p SkalType) CallArg {
	arg := NewCallArg(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Values
		case token.Value:
			if len(child.Children) == 0 {
				continue
			}
			value := buildValue(child, p)
			arg.Values = append(arg.Values, &value)

		// '...'
		case token.Spread:
			arg.Spread = true

		default:
			sklog.UnexpectedType("typeset call arg node", child.Type.String())
		}
	}

	return arg
}
