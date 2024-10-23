package validate

import (
	"fmt"
	"slices"
	"strings"

	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/lua"
	"github.com/illbjorn/skal/internal/skal/sklog"
	"github.com/illbjorn/skal/internal/skal/typeset"
	"github.com/illbjorn/skal/pkg/fstr"
)

var sprintf = fmt.Sprintf

func Validate(set *typeset.TypeSet) {
	for _, member := range set.Members {
		switch t := member.Value.(type) {
		case *typeset.Struct:
			validateStruct(t)

		case *typeset.Enum:
			validateEnum(t)

		case *typeset.Bind:
			validateBind(t)

		case *typeset.Statement:
			validateStatement(t, nil)

		case *typeset.Fn:
			validateFn(t)

		case *typeset.Call:
			validateCall(t)

		case []*typeset.External:
			validateExternal(t)

		default:
			// TODO
			// println("Found:", sprintf("%T", t))
		}
	}
}

func validateExternal(ext []*typeset.External) {
	// Register
	for _, e := range ext {
		v.Register(e.Alias, token.Fn, token.Fn, "", e)
	}
}

func validateStruct(st *typeset.Struct) {
	// Register Struct
	v.Register(st.Ref(), token.Struct, token.Struct, "", st)

	// Register Fields
	for _, field := range st.Fields {
		v.Register(field.Ref(), token.StructField, token.StructField, st.Ref(), field)
	}

	// Register methods, validate method blocks.
	for _, method := range st.Methods {
		// Register method name.
		v.Register(method.Ref(), token.Fn, token.StructMethod, st.Ref(), method)

		// Register method args.
		for _, arg := range method.Args {
			v.Register(arg.Ref(), token.FnArg, token.FnArg, st.Ref(), arg)
		}

		// Validate all block statements.
		for _, stmt := range method.Block {
			validateStatement(stmt, st)
		}
	}
}

func validateEnum(enum *typeset.Enum) {
	v.Register(enum.Ref(), token.Enum, token.Enum, "", enum)

	for _, member := range enum.Members {
		v.Register(member.Ref(), member.ValueType, token.EnumMember, enum.Ref(), member)

		// All enum members must match their parent's value type.
		if member.ValueType != enum.MemberType {
			err := "Enum member values must all have the same type as the parent enum."
			verr(err, member.Token())
		}
	}
}

func validateBind(bind *typeset.Bind) {
	if !bind.Rebind {
		for _, b := range bind.Binds {
			v.Register(b.Ref(), token.Bind, token.Bind, "", bind)
		}
		return
	}

	// Validate we're rebinding a value that exists.
	for _, b := range bind.Binds {
		validateRef(b, bind)
	}
}

func validateFn(fn *typeset.Fn) {
	// Register fn.
	v.Register(fn.Ref(), token.Fn, token.Fn, "", fn)

	// Register args.
	for _, arg := range fn.Args {
		v.Register(arg.Ref(), token.FnArg, token.FnArg, fn.Ref(), arg)
	}

	// Validate statements.
	for _, stmt := range fn.Block {
		validateStatement(stmt, fn)
	}
}

func validateCall(call *typeset.Call) {
	// The fn definition must exist.
	validateRef(call.SkalType, call)

	kt := v.Select().ID(call.Ref()).SkalType(token.Fn).First()
	if kt == nil {
		err := fstr.Pairs(
			"Failed to locate called fn: {fn}.",
			"fn", call.Ref(),
		)
		verr(err, call.Token())
	}

	fn := kt.object.(*typeset.Fn)

	// Confirm the number of args match.
	if len(call.Args) != len(fn.Args) {

	}

	// Validate arg types match.
	for _, callArg := range call.Args {
		for _, fnArg := range fn.Args {
			callArgType := getValueType(callArg.Values)
			if fnArg.Type() != callArgType {
				err := fstr.Pairs(
					"Cannot use value of type '{valueT}' for '{argT}' arg '{argN}' in call to '{fn}'.",
					"valueT", callArgType,
					"argT", fnArg.Type(),
					"argN", fnArg.ID(),
					"fn", fn.Ref(),
				)

				verr(err, call.Token())
			}
		}
	}
}

func getValueType(values []*typeset.Value) string {
	t := ""
	for _, value := range values {
		if value.Type() != token.Undefined {
			if t == "" {
				t = value.Type()
			}

			if value.Type() != t {
				println("Uh, oh.")
			}

			continue
		}

		switch value.ValueType {
		case token.Ref: // TODO: Lookup the ref and get its type.
			t := lookupRefType(value.Ref(), v)

			println(t)
		default:
			sklog.UnexpectedType("get value", value.ValueType)
		}
	}

	return t
}

func validateStatement(stmt *typeset.Statement, parent typeset.SkalType) {
	switch stmt.StmtType {
	case token.If:
		for _, cond := range stmt.If.Conditions {
			validateValue(cond, parent)
		}

	case token.Bind, token.Rebind:
		validateBind(stmt.Bind)

	case token.Call:
		validateCall(stmt.Call)

	case token.Ret:
		for _, v := range stmt.Values {
			validateValue(v, parent)
		}

	case token.For:

	default:
		sklog.UnexpectedType("validate statement", stmt.StmtType)
	}
}

// TODO: Following the type system implementation, we need to validate all of
// the below are compatible with preceding / following values.
// FEATURE-BLOCK: Types.
func validateValue(value *typeset.Value, parent typeset.SkalType) {
	switch value.ValueType {
	case token.Ref:
		validateRef(value.SkalType, parent)

	case token.Call:
		validateCall(value.Call)

	case
		token.ComparisonOperator, token.IntL, token.StrL, token.BoolL, token.Not,
		token.Nil, token.LogicOperator, token.MathOperator,
		token.ConcatOperator:

	case token.ValueGroup:
		for _, v := range value.Group {
			validateValue(v, parent)
		}

	default:
		sklog.UnexpectedType("validate value value", value.ValueType)
	}
}

// TODO: Validate positive ref matches have an expected `SkalT` (e.g. `Method`).
// TODO: Target fields or methods more specifically when looking up a `this` ref.
// TODO: Fix the wonky behavior with index refs including the `[]`.
func validateRef(base typeset.SkalType, parent typeset.SkalType) {
	refs := base.Refs()

	// Ignore stdlib fns.
	if slices.Contains(lua.StdlibFns, base.Ref()) {
		return
	}

	// var baseT string
	var i int
	var ref string
outer:
	for {
		// Break when we hit the length.
		if i == len(refs) {
			return
		}

		ref = refs[i]

		// SEE TODOs.
		if strings.HasPrefix(ref, "[") {
			ref = strings.ReplaceAll(ref, "[", "")
			ref = strings.ReplaceAll(ref, "]", "")
		}

		// Prepare a token, referenced if we hit an error.
		tk := base.Token()
		if tk == nil {
			tk = parent.Token()
		}

		//--------------------------------------------------------------------------
		// 'this'
		//
		// If we find a `this` reference, confirm the parent is a struct and the
		// struct has a field matching refs[i+1].
		if ref == token.This {
			res := ascendTo[*typeset.Struct](parent)

			// Get the field name.
			// If this is the last ref in the slice, it's likely just a return
			// returning the instance.
			if i+1 >= len(refs) {
				return
			}

			fieldName := refs[i+1]

			// First check fields.
			// SEE TODOs.
			for _, field := range res.Fields {
				if fieldName == field.Ref() {
					i = i + 2
					continue outer
				}
			}

			switch parent.(type) {
			// Call
			case *typeset.Call:
				for _, method := range res.Methods {
					if fieldName == method.ID() {
						i = i + 2
						continue outer
					}
				}

			default:
				err := fstr.Pairs(
					"Validate ref, `this` handling, found parent type: {type}.",
					"type", sprintf("%T", parent),
				)
				verr(err, tk)
				return
			}

			// If we made it here, the ref doesn't exist.
			err := fstr.Pairs(
				"Failed to locate field or method: '{field}' on parent: '{parent}'.",
				"field", fieldName,
				"parent", res.Ref(),
			)
			verr(err, tk)
			return
		}

		//--------------------------------------------------------------------------
		// All other refs.
		//
		// Lookup the ref, confirm it exists.
		res := v.Select().ID(ref).First()
		if res == nil {
			verr("Found undefined reference: "+ref+".", tk)
			return
		}

		// Increment the index.
		i++
	}
}
