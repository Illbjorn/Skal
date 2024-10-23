package typeset

import (
	"bytes"
	"strings"

	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/lua"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

type SkalType interface {
	Type() string
	SetType(string)
	Pub() bool
	SetPub()
	ID() string
	AddRef(r string)
	Ref() string
	Refs() []string
	MethodRef() string
	RefsLen() int
	SetParent(p SkalType)
	Parent() SkalType
	AddDefer(d *Statement)
	Defers() chan *Statement
	SetToken(tk token.Token)
	Token() token.Token
}

var _ SkalType = new(Base)

func NewBase(n *parse.Node, p SkalType) *Base {
	return &Base{
		parent: p,
		token:  n.Token,
		_type:  token.Undefined,
	}
}

type Base struct {
	_type  string
	parent SkalType
	// token  token.Token
	token  token.Token
	refs   []string
	defers []*Statement
	pub    bool
}

func (s *Base) Type() string     { return s._type }
func (s *Base) SetType(t string) { s._type = t }
func (s *Base) Pub() bool        { return s.pub }
func (s *Base) SetPub()          { s.pub = true }

func (s *Base) ID() string {
	if s.RefsLen() == 0 {
		return ""
	}

	return s.refs[s.RefsLen()-1]
}

func (s *Base) AddRef(r string) {
	s.refs = append(s.refs, r)
}

func (s *Base) Ref() string {
	if s.RefsLen() == 0 {
		return ""
	}

	if s.RefsLen() == 1 {
		return s.refs[0]
	}

	out := bytes.NewBuffer(nil)
	var next string
	for i, p := range s.refs {
		if i < s.RefsLen()-1 {
			next = s.refs[i+1]
		} else {
			next = ""
		}

		if p == token.Comma {
			out.WriteString(", ")
			continue
		}

		out.WriteString(lua.Translate(p))
		if i < s.RefsLen()-1 && !strings.HasPrefix(
			next,
			"[") && next != token.Comma {
			out.WriteString(".")
		}
	}

	return out.String()
}

func (s *Base) Refs() []string {
	return s.refs
}

func (s *Base) MethodRef() string {
	// Skip the trouble of allocating.
	if s.RefsLen() == 0 {
		return ""
	}

	if s.RefsLen() == 1 {
		return lua.Translate(s.refs[0])
	}

	if s.RefsLen() == 2 {
		return lua.Translate(s.refs[0]) + ":" + s.refs[1]
	}

	out := bytes.NewBuffer(nil)
	var next string
	for i, p := range s.refs {
		if i < s.RefsLen()-1 {
			next = s.refs[i+1]
		} else {
			next = ""
		}

		if p == token.Comma {
			out.WriteString(", ")
			continue
		}

		out.WriteString(lua.Translate(p))
		if i < s.RefsLen()-2 && !strings.HasPrefix(
			next,
			"[") && next != token.Comma {
			out.WriteString(".")
		}

		if i == s.RefsLen()-2 {
			out.WriteString(":")
		}
	}

	return out.String()
}

func (s *Base) RefsLen() int          { return len(s.refs) }
func (s *Base) SetParent(p SkalType)  { s.parent = p }
func (s *Base) Parent() SkalType      { return s.parent }
func (s *Base) AddDefer(d *Statement) { s.defers = append(s.defers, d) }
func (s *Base) Defers() chan *Statement {
	ch := make(chan *Statement)

	go func() {
		defer close(ch)
		for _, d := range s.defers {
			ch <- d
		}
	}()

	return ch
}

func (s *Base) SetToken(tk token.Token) {
	if tk == nil {
		return
	}

	s.token = tk
}

func (s *Base) Token() token.Token {
	if s.token == nil {
		return nil
	}

	return s.token
}

/*------------------------------------------------------------------------------
 * Ref Build
 *----------------------------------------------------------------------------*/

func buildRef[T SkalType](n node, t T) T {
	for _, child := range n.Children {
		switch child.Type {
		// 'this' | ID
		case token.This, token.ID:
			t.AddRef(child.Value)

		// Index
		// [0]
		// [this.Nested.Property]
		// [this.Call()]
		case token.Index:
			t.AddRef("[" + strings.Join(buildRefIndex(child), ".") + "]")

		default:
			sklog.UnexpectedType("typeset ref node", child.Type)
		}
	}

	return t
}

func buildRefIndex(n node) []string {
	var out []string
	for _, child := range n.Children {
		if child.Value == token.This {
			// 'this'
			out = append(out, "self")
		} else {
			// Anything else.
			out = append(out, child.Value)
		}
	}
	return out
}
