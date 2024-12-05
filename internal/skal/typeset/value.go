package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewValue(n *parse.Node, p SkalType) Value {
	return Value{SkalType: NewBase(n, p)}
}

type Value struct {
	SkalType
	Call       *Call
	Fn         *Fn
	IntL       string
	BoolL      string
	StrL       string
	Nil        string
	Op         string
	List       []*Value
	Group      []*Value
	Comparison []*Value
	ValueType  token.Type
	Not        bool
}

func buildValue(n node, p SkalType) Value {
	v := NewValue(n, p)
	v.ValueType = n.Type

	switch n.Type {
	case token.Value:
		return buildValue(n.Children[0], p)

	// Value group.
	case token.ValueGroup:
		return buildValueGroup(n, p)

	// 'fn'
	case token.Fn:
		v.ValueType = token.Fn
		fn := buildFn(n, p)
		v.Fn = &fn

	// Reference
	case token.Ref:
		v = *buildRef(n, &v)

	// Call
	case token.Call:
		call := buildCall(n, p)
		v.Call = &call

	// BoolL
	case token.BoolL:
		v.SetType(token.Bool)
		v.BoolL = n.Value

	// StrL
	case token.StrL:
		v.SetType(token.Str)
		v.StrL = n.Value

	// IntL
	case token.IntL:
		v.SetType(token.Int)
		v.IntL = n.Value

	// Nil
	case token.Nil:
		v.SetType(token.Nil)
		v.Nil = n.Value

	// List
	case token.ListL:
		v.SetType(token.List)
		for _, child := range n.Children {
			value := buildValue(child, p)
			v.List = append(v.List, &value)
		}

	// '[]'
	case token.List: // Nothing to do.
		v.SetType(token.List)

	// '!'
	case token.Not:
		v.Not = true

	// All Operators
	case token.ConcatOperator, token.MathOperator, token.ComparisonOperator, token.LogicOperator:
		v.Op = n.Value

	default:
		sklog.UnexpectedType("typeset value node", n.Type.String())
	}

	return v
}

func buildValueGroup(n node, p SkalType) Value {
	group := NewValue(n, p)
	group.ValueType = token.ValueGroup

	for _, child := range n.Children {
		value := buildValue(child, p)
		group.Group = append(group.Group, &value)
	}

	return group
}
