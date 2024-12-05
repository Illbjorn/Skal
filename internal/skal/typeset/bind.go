package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewBind(n *parse.Node, p SkalType) Bind {
	return Bind{SkalType: NewBase(n, p)}
}

type Bind struct {
	SkalType
	ValueType string
	Binds     []*Base
	Values    []*Value
	Rebind    bool
}

func buildBind(n node, p SkalType, rebind bool) Bind {
	bind := NewBind(n, p)
	bind.Rebind = rebind

	for _, child := range n.Children {
		switch child.Type {
		// Reference
		case token.Ref:
			bind.Binds = append(bind.Binds, buildRef(child, &Base{}))

		// Values
		case token.Value:
			value := buildValue(child, &bind)
			bind.Values = append(bind.Values, &value)

		default:
			sklog.UnexpectedType("typeset bind node", child.Type.String())
		}
	}

	return bind
}
