package typeset

import "github.com/illbjorn/skal/internal/skal/lex/token"

func NewTypeSet() TypeSet {
	return TypeSet{
		Members: make([]Type, 0),
	}
}

type TypeSet struct {
	Members []Type
}

// Add adds a given set of instance information to the members slice.
func (c *TypeSet) Add(v any, id string, t token.Type) {
	if v == nil {
		return
	}

	// Append the object to the members slice.
	c.Members = append(c.Members, Type{ID: id, Value: v, Type: t})
}

type Type struct {
	Value any
	ID    string
	Type  token.Type
}
