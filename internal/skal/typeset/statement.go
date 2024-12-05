package typeset

import (
	"fmt"

	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

/*------------------------------------------------------------------------------
 * Statement
 *----------------------------------------------------------------------------*/

func NewStatement(n *parse.Node, p SkalType) Statement {
	return Statement{SkalType: NewBase(n, p)}
}

type Statement struct {
	SkalType
	If       *If
	For      *For
	Call     *Call
	Fn       *Fn
	Bind     *Bind
	Op       string
	Values   []*Value
	StmtType token.Type
}

func buildStatement(n node, p SkalType) Statement {
	stmt := NewStatement(n, p)

	for _, child := range n.Children {
		switch child.Type {
		// 'return'
		case token.Ret:
			stmt.StmtType = token.Ret
			for _, child := range child.Children {
				value := buildValue(child, p)
				stmt.Values = append(stmt.Values, &value)
			}

		// 'for'
		case token.For:
			stmt.StmtType = token.For
			nfor := buildFor(child, &stmt)
			stmt.For = &nfor

		// Call
		case token.Call:
			stmt.StmtType = token.Call
			call := buildCall(child, &stmt)
			stmt.Call = &call

		// 'if'
		case token.If:
			stmt.StmtType = token.If
			nif := buildIf(child, &stmt)
			stmt.If = &nif

		// 'let'
		case token.Bind:
			stmt.StmtType = token.Bind
			bind := buildBind(child, &stmt, false)
			stmt.Bind = &bind

		// Rebind
		case token.Rebind:
			stmt.StmtType = token.Rebind
			bind := buildBind(child, &stmt, true)
			bind.Rebind = true
			stmt.Bind = &bind

		// 'fn'
		case token.Fn:
			stmt.StmtType = token.Fn
			fn := buildFn(child, &stmt)
			stmt.Fn = &fn

		// 'defer'
		case token.Defer:
			def := buildStatement(child.Children[0], p)
			stmt.StmtType = token.Defer
			stmt.AddDefer(&def)

		case token.Ref:
			fmt.Println(child.Children)

		default:
			sklog.UnexpectedType("typeset statement node", child.Type.String())
		}
	}

	if n.Token != nil {
		fmt.Println(n.Token.File())
	}

	if stmt.StmtType == 0 {
		sklog.CFatal("Statement fell through with no type.")
	}

	return stmt
}

func buildBlock(n node, p SkalType) []*Statement {
	var block []*Statement

	for _, child := range n.Children {
		switch child.Type {
		// Statement
		case token.Statement:
			stmt := buildStatement(child, p)
			block = append(block, &stmt)

		default:
			sklog.UnexpectedType("typeset block node", child.Type.String())
		}
	}

	return block
}
