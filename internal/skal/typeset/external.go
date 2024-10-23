package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewExternal(n *parse.Node, p SkalType) External {
	return External{SkalType: NewBase(n, p)}
}

type External struct {
	SkalType
	Alias  string
	ExtRef []string
}

func buildExtern(n node) []*External {
	extern := make([]*External, 0)

	for _, child := range n.Children {
		switch child.Type {
		// External
		case token.External:
			external := buildExternal(child, nil)
			extern = append(extern, &external)

		default:
			sklog.UnexpectedType("typeset extern", child.Type)
		}
	}

	return extern
}

func buildExternal(n node, p SkalType) External {
	ext := NewExternal(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			ext.AddRef(child.Value)

		// 'as'
		case token.As:
			ext.Alias = child.Value

		default:
			sklog.UnexpectedType("external node", child.Type)
		}
	}

	return ext
}
