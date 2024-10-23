package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewEnum(n *parse.Node, p SkalType) Enum {
	return Enum{SkalType: NewBase(n, p)}
}

type Enum struct {
	SkalType
	MemberType string
	Members    []*EnumMember
}

func NewEnumMember(n *parse.Node, p SkalType) EnumMember {
	return EnumMember{SkalType: NewBase(n, p)}
}

type EnumMember struct {
	SkalType
	Value     string
	ValueType string
}

func buildEnum(n node) Enum {
	e := NewEnum(n, nil)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			e.AddRef(child.Value)

		// Members
		case token.EnumMember:
			member := buildEnumMember(child, &e)
			e.Members = append(e.Members, &member)
			// TEMPORARY: We will eventually support enum typing in the syntax.
			e.MemberType = e.Members[0].ValueType

		default:
			sklog.UnexpectedType("typeset enum", child.Type)
		}
	}

	return e
}

func buildEnumMember(n node, p SkalType) EnumMember {
	m := NewEnumMember(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			m.AddRef(child.Value)

		// String Literal | Int Literal | Bool Literal
		case token.StrL, token.IntL, token.BoolL:
			m.ValueType = child.Type
			m.Value = child.Value

		default:
			sklog.UnexpectedType("typeset enum member", child.Type)
		}
	}

	return m
}
