package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func NewStruct(n *parse.Node, p SkalType) Struct {
	return Struct{SkalType: NewBase(n, p)}
}

type Struct struct {
	SkalType
	Fields        []*StructField
	Methods       []*Fn
	NoConstructor bool
}

func NewStructField(n *parse.Node, p SkalType) StructField {
	return StructField{SkalType: NewBase(n, p)}
}

type StructField struct {
	SkalType
}

func buildStruct(n node) Struct {
	s := NewStruct(n, nil)

	for _, child := range n.Children {
		switch child.Type {
		// ID
		case token.ID:
			s.AddRef(child.Value)

		// Field
		case token.StructField:
			nf := buildStructField(child, &s)
			s.Fields = append(s.Fields, &nf)

		// Method
		case token.StructMethod:
			fn := buildFn(child, &s)
			s.Methods = append(s.Methods, &fn)

		default:
			sklog.UnexpectedType("typeset struct node", child.Type.String())
		}
	}

	return s
}

func buildStructField(
	n node,
	p SkalType,
) StructField {
	f := NewStructField(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// Reference
		case token.Ref:
			f = *buildRef(child, &f)

		default:
			sklog.UnexpectedType("struct field node", child.Type.String())
		}
	}

	return f
}
