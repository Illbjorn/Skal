package parse

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func Parse(tc *token.Collection, root, n *Node) *Node {
	if root == nil {
		root = new(Node)
	}

	if n == nil {
		n = new(Node)
	}

	tk := tc.LA()

	switch tk.Type() {
	// 'pub'
	case token.Pub:
		n.AddChild(
			new(Node).SetToken(tc.Adv()),
		)
		return Parse(tc, root, n)

	// 'enum'
	case token.Enum:
		n.AddChild(
			parseEnum(tc),
		)

	// 'struct'
	case token.Struct:
		n.AddChild(
			parseStruct(tc),
		)

	// 'fn'
	case token.Fn:
		n.AddChild(
			parseFn(tc),
		)

	// 'if'
	case token.If:
		n.AddChild(
			parseIf(tc),
		)

	// 'if'
	case token.For:
		n.AddChild(
			parseFor(tc),
		)

	// 'let'
	case token.Let:
		n.AddChild(
			parseBind(tc),
		)

	// 'import'
	case token.Import:
		parseError(
			"Import statements must appear in the File before any other code.",
			tk,
			true,
		)

	// ID
	// Rebind | Call
	case token.ID:
		tk := tc.LookPastRef()
		switch tk.Type() {
		// '('
		// Call
		case token.ParenOpen:
			n.AddChild(
				parseCall(tc),
			)

		// '='
		// Rebind
		case token.EQ:
			n.AddChild(
				parseRebind(tc),
			)

		default:
			sklog.UnexpectedType("LookPast token", tk.Type())
		}

	// 'extern'
	case token.Extern:
		n.AddChild(
			parseExtern(tc),
		)

	// EOF
	case token.EOF:
		return root

	// 'defer'
	case token.Defer:
		parseError(
			"Top level deferrals are not allowed.",
			tk,
			true)

	default:
		sklog.UnexpectedType("parse token", tc.LA().Type())
	}

	root.AddChild(n)
	return Parse(tc, root, new(Node))
}

func parseDefer(tc *token.Collection) *Node {
	ndefer := new(Node).SetType(token.Defer).SetTokenOnly(tc.AdvT(token.Defer)).AddChild(
		parseStatement(tc, nil),
	)

	return ndefer
}

func parseExtern(tc *token.Collection) *Node {
	extern := new(Node).SetType(token.Extern).SetTokenOnly(tc.AdvT(token.Extern))

	// '{'
	tc.AdvT(token.BraceOpen)

	// Reference
	for !tc.NTT(token.BraceClose) {
		external := new(Node).SetType(token.External).SetTokenOnly(tc.LA())

		// Reference
		for {
			// ID
			external.AddChild(
				new(Node).SetToken(tc.AdvT(token.ID)),
			)

			// '.'
			if _, ok := tc.AdvIf(token.Dot); !ok {
				break
			}
		}

		// 'as'
		tc.AdvT(token.As)

		// Alias
		external.AddChild(
			new(Node).SetType(token.As).SetToken(tc.AdvT(token.ID)),
		)

		extern.AddChild(
			external,
		)
	}

	// '}'
	tc.AdvT(token.BraceClose)

	return extern
}

func parseStruct(tc *token.Collection) *Node {
	nstruct := new(Node).SetType(token.Struct).SetTokenOnly(tc.LA())

	// 'struct'
	tc.AdvT(token.Struct)

	// ID
	nstruct.AddChild(
		new(Node).SetToken(tc.AdvT(token.ID)),
	)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Fields
	nstruct.AddChildren(
		parseStructFields(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return nstruct
}

func parseStructFields(tc *token.Collection) []*Node {
	var fields []*Node

	for {
		if tc.LA().Type() == token.BraceClose {
			return fields
		}

		switch {
		// Method
		case tc.LookPastRef().Type() == token.ParenOpen:
			fields = append(fields, parseMethod(tc))

		// Method
		case tc.LA().Type() == token.New:
			fields = append(fields, parseMethod(tc))

		// Reference
		default:
			fields = append(
				fields,
				new(Node).SetType(token.StructField).SetTokenOnly(tc.LA()).AddChild(parseRef(tc)),
			)
		}
	}
}

func parseMethod(tc *token.Collection) *Node {
	fn := new(Node).SetType(token.StructMethod).SetTokenOnly(tc.LA())

	// ID | 'new'
	fn.AddChild(
		new(Node).SetToken(tc.AdvOneOfT(token.ID, token.New)),
	)

	// '('
	tc.AdvT(token.ParenOpen)

	// Args
	fn.AddChildren(
		parseFnArgs(tc),
	)

	// ')'
	tc.AdvT(token.ParenClose)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Block
	fn.AddChild(
		parseBlock(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return fn
}

func parseEnum(tc *token.Collection) *Node {
	enum := new(Node).SetType(token.Enum).SetTokenOnly(tc.LA())

	// 'enum'
	tc.AdvT(token.Enum)

	// ID
	enum.AddChild(
		new(Node).SetToken(tc.AdvT(token.ID)),
	)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Members
	enum.AddChildren(
		parseEnumMembers(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return enum
}

func parseEnumMembers(tc *token.Collection) []*Node {
	var members []*Node

	for {
		if tc.NTT(token.BraceClose) {
			return members
		}

		member := new(Node).SetType(token.EnumMember).SetTokenOnly(tc.LA())

		// ID
		member.AddChild(
			new(Node).SetToken(tc.AdvT(token.ID)),
		)

		// '='
		tc.AdvT(token.EQ)

		// Value
		member.AddChild(
			new(Node).SetToken(
				tc.AdvOneOfT(
					token.IntL,
					token.BoolL,
					token.StrL)),
		)

		members = append(members, member)
	}
}

func parseFn(tc *token.Collection) *Node {
	fn := new(Node).SetType(token.Fn).SetTokenOnly(tc.LA())

	// 'fn'
	tc.AdvT(token.Fn)

	// ID
	fn.AddChild(
		new(Node).SetToken(tc.AdvT(token.ID)),
	)

	// '('
	tc.AdvT(token.ParenOpen)

	// Args
	fn.AddChildren(
		parseFnArgs(tc),
	)

	// ')'
	tc.AdvT(token.ParenClose)

	// OPTIONAL: Return type hint.
	if tk, ok := tc.AdvIf(token.Fn, token.Str, token.Int, token.Bool, token.ID); ok {
		fn.AddChild(
			new(Node).SetToken(tk).SetType(token.TypeHint),
		)
	}

	// '{'
	tc.AdvT(token.BraceOpen)

	// Block
	fn.AddChild(
		parseBlock(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return fn
}

func parseFnArgs(tc *token.Collection) []*Node {
	if tc.LA().Type() == token.ParenClose {
		return nil
	}
	var args []*Node

	// Args
	for {
		arg := new(Node).SetType(token.FnArg).SetTokenOnly(tc.LA())

		// OPTIONAL: '...'
		if tk, ok := tc.AdvIf(token.Spread); ok {
			arg.AddChild(
				new(Node).SetToken(tk),
			)
		}

		// ID
		arg.AddChild(
			new(Node).SetToken(tc.AdvT(token.ID)),
		)

		// OPTIONAL: Type hint.
		if _, ok := tc.AdvIf(token.Colon); ok {
			// Type token.
			tk := tc.AdvOneOfT(
				token.ID,
				token.Str,
				token.Int,
				token.Bool,
				token.Fn,
			)

			// Add the hint.
			arg.AddChild(
				new(Node).SetToken(tk).SetType(token.TypeHint),
			)
		}

		args = append(args, arg)
		// OPTIONAL: ','
		if _, ok := tc.AdvIf(token.Comma); !ok {
			return args
		}
	}
}

func parseAnonFn(tc *token.Collection) *Node {
	fn := new(Node).SetType(token.Fn).SetTokenOnly(tc.LA())

	// '('
	tc.AdvT(token.ParenOpen)

	// OPTIONAL: Args
	if !tc.NTT(token.ParenClose) {
		var args []*Node

		// Arg
		for {
			arg := new(Node).SetType(token.FnArg).SetTokenOnly(tc.LA())

			// ID
			arg.AddChild(
				new(Node).SetToken(tc.AdvT(token.ID)),
			)

			args = append(args, arg)

			// OPTIONAL: ','
			if _, ok := tc.AdvIf(token.Comma); !ok {
				break
			}
		}

		fn.AddChildren(
			args,
		)
	}

	// ')'
	tc.AdvT(token.ParenClose)

	// '->'
	tc.AdvT(token.Arrow)

	// Value
	fn.AddChildren(
		parseValue(tc),
	)

	return fn
}

func parseReturnStatement(tc *token.Collection) *Node {
	ret := new(Node).SetToken(tc.AdvT(token.Ret))

	// Bare return.
	if tc.NTT(token.BraceClose) {
		return ret
	}

	// Returned Value.
	ret.AddChildren(
		parseValue(tc),
	)

	return ret
}

func parseIf(tc *token.Collection) *Node {
	nif := new(Node).SetType(token.If).SetTokenOnly(tc.LA())

	// 'if'
	tc.AdvT(token.If)

	// Conditions
	nif.AddChild(
		parseConditions(tc),
	)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Block
	nif.AddChild(
		parseBlock(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	// Elif+
	nif.AddChildren(
		parseElifs(tc),
	)

	// Else
	nif.AddChild(
		parseElse(tc),
	)

	return nif
}

func parseElifs(tc *token.Collection) []*Node {
	if tc.LA().Type() != token.Elif {
		return nil
	}
	var elifs []*Node

	for {
		elif := new(Node).SetType(token.Elif).SetTokenOnly(tc.LA())

		// 'elif'
		tc.AdvT(token.Elif)

		// Conditions.
		elif.AddChild(
			parseConditions(tc),
		)

		// '{'
		tc.AdvT(token.BraceOpen)

		// Block
		elif.AddChild(
			parseBlock(tc),
		)

		// '}'
		tc.AdvT(token.BraceClose)

		elifs = append(elifs, elif)
		if tc.LA().Type() != token.Elif {
			return elifs
		}
	}
}

func parseElse(tc *token.Collection) *Node {
	if tc.LA().Type() != token.Else {
		return nil
	}

	nelse := new(Node).SetType(token.Else)

	// 'else'
	tc.AdvT(token.Else)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Block
	nelse.AddChild(
		parseBlock(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return nelse
}

func parseConditions(tc *token.Collection) *Node {
	return new(Node).
		SetType(token.Conditions).
		SetTokenOnly(tc.LA()).
		AddChildren(parseValue(tc))
}

func parseBind(tc *token.Collection) *Node {
	bind := new(Node).SetType(token.Bind).SetTokenOnly(tc.LA())

	// 'let'
	tc.AdvT(token.Let)

	// Reference
	// Multiple binds can be performed in a single line, so we look for one or
	// more here.
	for {
		bind.AddChild(
			parseRef(tc),
		)

		if _, ok := tc.AdvIf(token.Comma); !ok {
			break
		}
	}

	// ';'
	// If we hit a semicolon, it's just bringing a variable into scope.
	if _, ok := tc.AdvIf(token.SemiColon); ok {
		return bind
	}

	// '='
	tc.AdvT(token.EQ)

	// Value
	bind.AddChildren(
		parseValue(tc),
	)

	return bind
}

func parseRebind(tc *token.Collection) *Node {
	reb := new(Node).SetType(token.Rebind).SetTokenOnly(tc.LA())

	// Reference
	reb.AddChild(
		parseRef(tc),
	)

	// '='
	tc.AdvT(token.EQ)

	// Value
	reb.AddChildren(
		parseValue(tc),
	)

	return reb
}

func parseFor(tc *token.Collection) *Node {
	nfor := new(Node).SetType(token.For).SetTokenOnly(tc.LA())

	// 'for'
	tc.AdvT(token.For)

	// Iterators
	nfor.AddChildren(
		parseForIterators(tc),
	)

	// '=' | 'in'
	nfor.AddChild(
		new(Node).SetType(token.ForType).SetToken(
			tc.AdvOneOfT(
				token.EQ,
				token.In)),
	)

	// Iterables
	nfor.AddChildren(
		parseForIterables(tc),
	)

	// '{'
	tc.AdvT(token.BraceOpen)

	// Block
	nfor.AddChild(
		parseBlock(tc),
	)

	// '}'
	tc.AdvT(token.BraceClose)

	return nfor
}

func parseBlock(tc *token.Collection) *Node {
	block := new(Node).SetType(token.Block).SetTokenOnly(tc.LA())

	for {
		if tc.NTT(token.BraceOpen, token.BraceClose) {
			return block
		}

		stmt := parseStatement(tc, nil)
		if stmt == nil {
			return block
		}

		block.AddChild(stmt)
	}
}

func parseForIterators(tc *token.Collection) []*Node {
	var iterators []*Node

	for {
		iterator := new(Node).SetType(token.ForIterator).SetTokenOnly(tc.LA())
		iterators = append(iterators, iterator)

		iterator.AddChild(
			parseRef(tc),
		)

		if _, ok := tc.AdvIf(token.Comma); !ok {
			return iterators
		}
	}
}

func parseForIterables(tc *token.Collection) []*Node {
	var iterables []*Node

	for {
		iterable := new(Node).SetType(token.ForIterable).SetTokenOnly(tc.LA())

		tk := tc.LA()
		switch tk.Type() {
		// Int Literal
		case token.IntL, token.StrL:
			iterable.AddChild(
				new(Node).SetToken(tc.Adv()),
			)

		// Ref
		case token.ID, token.This:
			iterable.AddChild(
				parseRef(tc),
			)

		default:
			sklog.UnexpectedType("for iterables token", tk.Type())
		}

		iterables = append(iterables, iterable)
		if _, ok := tc.AdvIf(token.Comma); !ok {
			return iterables
		}
	}
}

func parseStatement(tc *token.Collection, stmt *Node) *Node {
	if stmt == nil {
		stmt = new(Node).SetType(token.Statement)
	}

	tk := tc.LA()
	switch tk.Type() {
	// 'pub'
	case token.Pub:
		pub := new(Node).SetTokenOnly(tc.Adv())
		stmt.AddChild(pub)
		return parseStatement(tc, stmt)

	// Rebind | Call | Reference
	case token.ID, token.This:
		// Rebind
		if tc.LookPastRef().Type() == token.EQ {
			stmt.AddChild(
				parseRebind(tc),
			)
			return stmt
		}

		// Call
		if tc.LookPastRef().Type() == token.ParenOpen {
			stmt.AddChild(
				parseCall(tc),
			)
			return stmt
		}

		// Reference
		stmt.AddChild(
			parseRef(tc),
		)

	// 'if'
	case token.If:
		stmt.AddChild(
			parseIf(tc),
		)

	// 'for'
	case token.For:
		stmt.AddChild(
			parseFor(tc),
		)

	// 'let'
	case token.Let:
		stmt.AddChild(
			parseBind(tc),
		)

	// 'fn'
	case token.Fn:
		stmt.AddChild(
			parseFn(tc),
		)

	// 'defer'
	case token.Defer:
		stmt.AddChild(
			parseDefer(tc),
		)

	// 'return'
	case token.Ret:
		stmt.AddChild(
			parseReturnStatement(tc),
		)

	default:
		sklog.UnexpectedType("parse statement", tk.Type())
	}

	return stmt
}

func parseCall(tc *token.Collection) *Node {
	call := new(Node).SetType(token.Call).SetTokenOnly(tc.LA())

	// Ref
	call.AddChild(
		parseRef(tc),
	)

	// Args
	call.AddChildren(
		parseCallArgs(tc),
	)

	return call
}

func parseCallArgs(tc *token.Collection) []*Node {
	var args []*Node

	// '('
	tc.AdvT(token.ParenOpen)

	// Args
	for {
		arg := new(Node).SetType(token.CallArg).SetTokenOnly(tc.LA())

		// Value
		arg.AddChildren(
			parseValue(tc),
		)

		// OPTIONAL: '...'
		if tk, ok := tc.AdvIf(token.Spread); ok {
			arg.AddChild(
				new(Node).SetToken(tk),
			)
		}

		// Add the node.
		args = append(args, arg)

		// ','
		if _, ok := tc.AdvIf(token.Comma); !ok {
			break
		}
	}

	// ')'
	tc.AdvT(token.ParenClose)

	return args
}

func parseRef(tc *token.Collection) *Node {
	ref := new(Node).SetType(token.Ref).SetTokenOnly(tc.LA())

	for {
		ref.AddChild(
			new(Node).SetToken(tc.AdvOneOfT(token.This, token.ID)),
		)

		if tc.NTT(token.BrackOpen) {
			ref.AddChild(
				parseIndex(tc),
			)
		}

		if _, ok := tc.AdvIf(token.Dot); !ok {
			return ref
		}
	}
}

func parseIndex(tc *token.Collection) *Node {
	index := new(Node).SetType(token.Index).SetTokenOnly(tc.LA())

	// '['
	tc.AdvT(token.BrackOpen)

	for {
		if tc.LA().Type() == token.BrackClose {
			break
		}

		// 'this' | IntL | StrL | BoolL | ID | '.'
		tk := tc.AdvOneOfT(
			token.This,
			token.IntL,
			token.StrL,
			token.BoolL,
			token.ID,
			token.Dot)

		// Drop dot operators.
		if tk.Type() == token.Dot {
			continue
		}

		index.AddChild(
			new(Node).SetToken(tk),
		)
	}

	// ']'
	tc.AdvT(token.BrackClose)

	return index
}

var ttsAnyOperator = []string{
	// String Concat
	token.Concat,
	// Math
	token.Mult, token.Div, token.Plus, token.Minus,
	// Comparison
	token.GE, token.LE, token.NE, token.EQEQ, token.GT, token.LT,
	// Logic
	token.And, token.Or, token.Not,
}

func parseValue(tc *token.Collection) []*Node {
	var values []*Node

	for {
		tk := tc.LA()
		value := new(Node).SetType(token.Value).SetTokenOnly(tk)

		switch tk.Type() {
		// '('
		// This can indicate an anonymous function or a group.
		case token.ParenOpen:
			if tc.LineAheadContains(token.Arrow) {
				value.AddChild(parseAnonFn(tc))

			} else {
				value.AddChild(parseValueGroup(tc))
			}

		// '!'
		case token.Not:
			value.AddChild(new(Node).SetToken(tc.Adv()))
			// This operator must be prefixed to a Value, so we restart the loop.
			values = append(values, value)
			continue

		// Call | Reference
		case token.ID, token.This:
			switch tc.LookPastRef().Type() {
			case token.ParenOpen:
				value.AddChild(parseCall(tc))

			default:
				value.AddChild(parseRef(tc))
			}

		// StrL
		case token.StrL:
			value.AddChild(new(Node).SetToken(tc.Adv()))

		// IntL
		case token.IntL:
			value.AddChild(new(Node).SetToken(tc.Adv()))

		// BoolL
		case token.True, token.False:
			value.AddChild(new(Node).SetToken(tc.Adv()).SetType(token.BoolL))

			// If there isn't a logic operator ahead, return.
			if !tc.NTT(token.And, token.Or) {
				values = append(values, value)
				return values
			}

		// Nil
		case token.Nil:
			value.AddChild(new(Node).SetToken(tc.Adv()))

		// '[]'
		case token.List:
			value.AddChild(
				new(Node).SetType(token.List).SetTokenOnly(tc.AdvT(token.List)),
			)

		// List Literal
		case token.BrackOpen:
			list := new(Node).SetTokenOnly(tc.LA()).SetType(token.ListL)

			// '['
			tc.AdvT(token.BrackOpen)

			for !tc.NTT(token.BrackClose) {
				switch tc.LA().Type() {
				// StrL | IntL
				case token.StrL, token.IntL:
					list.AddChild(
						new(Node).SetToken(tc.Adv()),
					)

				default:
					sklog.UnexpectedType(
						"parse list literal Value",
						tc.LA().Type())
				}

				// ','
				if _, ok := tc.AdvIf(token.Comma); !ok {
					break
				}
			}

			// ']'
			tc.AdvT(token.BrackClose)
			value.AddChild(list)

		default:
			// Value consumed.
			return values
		}

		values = append(values, value)

		// Operators
		// >= | <= | != | == | > | <
		// .. | +  | -  | *  | /
		// || | && | !
		if tk, ok := tc.AdvIf(ttsAnyOperator...); !ok {
			return values
		} else {
			opNode := new(Node).SetType(token.Value).SetTokenOnly(tc.LA())
			// Assign the operator to a broader group.
			ot := opType(tk.Type())
			opNode.AddChild(new(Node).SetType(ot).SetToken(tk))
			values = append(
				values,
				opNode,
			)
		}
	}
}

func parseValueGroup(tc *token.Collection) *Node {
	group := new(Node).SetType(token.ValueGroup).SetTokenOnly(tc.LA())
	// '('
	tc.AdvT(token.ParenOpen)

	// Value
	group.AddChildren(
		parseValue(tc),
	)

	// ')'
	tc.AdvT(token.ParenClose)

	return group
}

// Classifies a given operator under a broader set of groups for easier building
// and emitting later on.
func opType(op string) string {
	switch op {
	// '!'
	case token.Not:
		return token.Not

	// Comparison Operators
	case token.GE, token.LE, token.EQEQ, token.NE, token.GT, token.LT:
		return token.ComparisonOperator

	// Concat Operator
	case token.Concat:
		return token.ConcatOperator

	// Arithmetic Operators
	case token.Plus, token.Minus, token.Mult, token.Div:
		return token.MathOperator

	// Logic Operators
	case token.And, token.Or:
		return token.LogicOperator

	default:
		sklog.CFatalF(
			"Failed to classify binary operator: {op}.",
			"op", op,
		)
		return ""
	}
}
