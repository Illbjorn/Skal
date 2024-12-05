package typeset

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

// Staged as a local custom type to be eventually converted into an interface to
// detangle package coupling.
type node *parse.Node

func Typeset(tree node) TypeSet {
	ctc := NewTypeSet()

	for _, child := range tree.Children {
		ctc.Add(
			typeset(child),
		)
	}

	return ctc
}

func typeset(n node) (any, string, token.Type) {
	var pub bool
	for _, child := range n.Children {
		switch child.Type {
		// 'pub'
		case token.Pub:
			pub = true
			continue

		// Enum
		case token.Enum:
			enum := buildEnum(child)
			if pub {
				enum.SetPub()
			}
			return &enum, enum.Ref(), child.Type

		// Struct
		case token.Struct:
			nstruct := buildStruct(child)
			if pub {
				nstruct.SetPub()
			}
			return &nstruct, nstruct.Ref(), child.Type

		// Bind
		case token.Bind:
			bind := buildBind(child, nil, false)
			if pub {
				bind.SetPub()
			}
			return &bind, bind.Ref(), child.Type

		// Rebind
		case token.Rebind:
			bind := buildBind(child, nil, true)
			return &bind, bind.Ref(), child.Type

		// Fn
		case token.Fn:
			fn := buildFn(child, nil)
			if pub {
				fn.SetPub()
			}
			return &fn, fn.Ref(), token.Fn

		// 'if'
		case token.If:
			nif := buildIf(child, nil)
			return &nif, "", token.If

		// 'for'
		case token.For:
			nfor := buildFor(child, nil)
			return &nfor, "", token.For

		// Call
		case token.Call:
			call := buildCall(child, nil)
			return &call, call.Ref(), token.Call

		// 'extern'
		case token.Extern:
			extern := buildExtern(child)
			return extern, "", token.Extern

		default:
			sklog.UnexpectedType("typeset node", child.Type.String())
		}
	}

	return nil, "", 0
}
