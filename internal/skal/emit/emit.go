package emit

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/lua"
	"github.com/illbjorn/skal/internal/skal/sklog"
	"github.com/illbjorn/skal/internal/skal/typeset"
	"github.com/illbjorn/skal/pkg/formatter"
)

func Emit(
	ctc typeset.TypeSet,
	path string,
	isImport bool,
	f *formatter.Formatter,
) []byte {
	if len(ctc.Members) == 0 {
		return nil
	}

	// If we're processing an import, wrap it in a `do` block.
	// This allows us to enforce "cross-module" visibility boundaries.
	if isImport {
		// 'do'
		f.Newline().Str(stack.Indent()).Str("do")

		// Indicate the source File name as a line comment after the `do` block
		// open.
		baseName := filepath.Base(path)
		baseName = strings.Replace(baseName, filepath.Ext(baseName), "", -1)
		f.Str(" -- FILE: " + baseName)

		stack.Push()
	}

	for _, o := range ctc.Members {
		switch o.Type {
		// 'for'
		case token.For:
			s := o.Value.(*typeset.For)
			f.Newline().
				Str(emitFor(s))

		// 'enum'
		case token.Enum:
			s := o.Value.(*typeset.Enum)
			f.Newline().
				Str(emitEnum(s))

		// 'struct'
		case token.Struct:
			s := o.Value.(*typeset.Struct)
			f.Newline().
				Str(emitStruct(s))

		// Bind
		case token.Bind:
			s := o.Value.(*typeset.Bind)
			f.Newline().
				Str(emitBind(s))

		// 'fn'
		case token.Fn:
			s := o.Value.(*typeset.Fn)
			f.Newline().
				Str(emitFn(s))

		// Call
		case token.Call:
			s := o.Value.(*typeset.Call)
			f.Newline().
				Str(emitCall(s, false, false))

		// Rebind
		case token.Rebind:
			s := o.Value.(*typeset.Bind)
			f.Newline().
				Str(emitBind(s))

		// 'extern'
		case token.Extern:
			s := o.Value.([]*typeset.External)
			f.Str(emitExtern(s))

		// 'if'
		case token.If:
			s := o.Value.(*typeset.If)
			f.Newline().
				Str(emitConditional(s))

		default:
			sklog.UnexpectedType("emit member", o.Type)
		}
	}

	// If we're processing an import, close the open `do` block.
	if isImport {
		stack.Pop()
		f.Newline().
			Str(stack.Indent()).
			Str("end")
	}

	// Return the compiled code.
	return f.Bytes()
}

/*------------------------------------------------------------------------------
 * Extern
 *----------------------------------------------------------------------------*/

func emitExtern(ext []*typeset.External) string {
	f := formatter.NewFormatter()

	// For fancy output where each list member's assignment operator is aligned
	// we first process the externals, identifying the longest member and building
	// a slice of [`id`, `Value`] slices.
	var longest int
	exts := make([][2]string, len(ext))
	for i, s := range ext {
		if l := len(s.Alias); l > longest {
			longest = len(s.Alias)
		}
		exts[i] = [2]string{s.Alias, s.Ref()}
	}

	// Now we emit using each [2]string{}'s `0` index as the whitespace offset
	// between the ID and the `=` operator. Then we append the actual aliased
	// Value (`1` index) at the end.
	for _, ext := range exts {
		spaces := strings.Repeat(" ", longest-len(ext[0]))
		f.Newline().
			Str(stack.Indent()).
			// Alias
			Str(ext[0]).
			Str(spaces).
			// ' = '
			Str(" = ").
			// External
			Str(ext[1])
	}

	return f.String()
}

/*------------------------------------------------------------------------------
 * Structs
 *----------------------------------------------------------------------------*/

var tmplStruct = `
{in}{local}{ref} = \{}
{in}setmetatable({ref}, {ref})
{in}{ref}.__index = {ref}
{constructor}
`

var tmplStructNoCon = `
{in}{local}{ref} = \{}
{in}setmetatable({ref}, {ref})
{in}{ref}.__index = {ref}
`

func emitStruct(nstruct *typeset.Struct) string {
	f := formatter.NewFormatter()

	// Local
	var local string
	if !nstruct.Pub() {
		local = "local "
	}

	// Struct
	if !nstruct.NoConstructor {
		// Produce with constructor.
		f.Str(
			pairs(
				tmplStruct,
				"in", stack.Indent(),
				"local", local,
				"ref", nstruct.Ref(),
				"constructor", emitStructConstructor(nstruct),
			))
	} else {
		// Produce without constructor.
		f.Str(
			pairs(
				tmplStructNoCon,
				"in", stack.Indent(),
				"local", local,
				"ref", nstruct.Ref(),
			))
	}

	// Methods
	for _, method := range nstruct.Methods {
		f.Newline().
			Str(emitMethod(nstruct, method))
	}

	return f.String()
}

var tmplStructDefaultConstructor = `
{in}function {ref}:__call({args})
{instance}
{in}end
`

var tmplStructDefaultConstructorInstance = `
{in}return setmetatable(\{
{fields}
{in}}, self)
`

var tmplStructDefaultConstructorField = "{in}{ref} = {ref}"

func emitStructConstructor(nstruct *typeset.Struct) string {
	stack.Push() // <-- Constructor fn scope.

	// Assemble comma-delimited args to use as constructor function args.
	args := formatter.NewFormatter()
	// Assemble 'reference = reference' field pairs for default table assignments.
	fields := formatter.NewFormatter()
	stack.Push() // <-- Constructor fn instance member scope.
	for i, f := range nstruct.Fields {
		// Write the argument.
		args.Str(f.Ref())

		// Format and write the field initializer.
		fields.Str(
			pairs(
				tmplStructDefaultConstructorField,
				"in", stack.Indent(),
				"ref", f.Ref(),
			))

		if i < len(nstruct.Fields)-1 {
			args.Str(", ")
			fields.Str(",\n")
		}
	}
	stack.Pop() // Constructor fn instance member scope. --!>

	// Typeset the constructed instance.
	instance := pairs(
		tmplStructDefaultConstructorInstance,
		"in", stack.Indent(),
		"fields", fields.String(),
	)
	stack.Pop() // Constructor fn scope. --!>

	return pairs(
		tmplStructDefaultConstructor,
		"in", stack.Indent(),
		"ref", nstruct.Ref(),
		"args", args.String(),
		"instance", instance,
	)
}

var tmplMethod = `
{in}function {struct}:{fn}({args}){block}
{in}end
`

func emitMethod(nstruct *typeset.Struct, fn *typeset.Fn) string {
	// Args
	args := formatter.NewFormatter()
	var varArgName string
	for i, arg := range fn.Args {
		if arg.Vararg { // Vararg
			args.Str("...")
			// Indicate we need to table-capture the vararg following the fn open
			// definition.
			varArgName = arg.Ref()
		} else { // Arg
			args.Str(arg.Ref())
			if i < len(fn.Args)-1 {
				args.Str(", ")
			}
		}
	}

	// ID
	var id string
	if fn.Constructor {
		id = "__call"
	} else {
		id = fn.ID()
	}

	return pairs(
		tmplMethod,
		"in", stack.Indent(),
		"struct", nstruct.ID(),
		"fn", id,
		"args", args.String(),
		"block", emitFnBlock(fn, varArgName),
	)
}

/*------------------------------------------------------------------------------
 * Enums
 *----------------------------------------------------------------------------*/

var tmplEnum = `
{in}{local}{ref} = \{{members}
{in}}
`

func emitEnum(enum *typeset.Enum) string {
	f := formatter.NewFormatter()

	// Local
	var local string
	if !enum.Pub() {
		local = "local "
	}

	// Members
	members := enumMembers(enum.Members)

	return f.Str(
		pairs(
			tmplEnum,
			"local", local,
			"ref", enum.Ref(),
			"members", members,
			"in", stack.Indent(),
		)).String()
}

var tmplEnumMember = "{in}{id} = {Value}"

func enumMembers(members []*typeset.EnumMember) string {
	stack.Push()
	defer stack.Pop()

	// Members
	f := formatter.NewFormatter()
	for i, m := range members {
		value := m.Value
		if m.ValueType == token.StrL {
			value = "'" + value + "'"
		}

		f.Newline().Str(
			pairs(
				tmplEnumMember,
				"in", stack.Indent(),
				"id", m.Ref(),
				"Value", value,
			))

		if i < len(members)-1 {
			f.Str(",")
		}
	}

	return f.String()
}

/*------------------------------------------------------------------------------
 * Fns
 *----------------------------------------------------------------------------*/

var tmplFn = `
{in}{local}function {ref}({args}){block}
{in}end
`

func emitFn(fn *typeset.Fn) string {
	// Local
	var local string
	if !fn.Pub() && fn.RefsLen() <= 1 && !fn.Constructor {
		local = "local "
	}

	// Args
	args := formatter.NewFormatter()
	var varArgName string
	for i, arg := range fn.Args {
		if arg.Vararg { // Vararg
			// Add the vararg to our args output.
			args.Str("...")

			// Indicate we need to table-capture the vararg following the fn open
			// definition.
			varArgName = arg.Ref()
		} else { // Arg
			args.Str(arg.Ref())

			if i < len(fn.Args)-1 {
				args.Str(", ")
			}
		}
	}

	// ID
	var ref string
	if fn.RefsLen() > 1 {
		ref = fn.MethodRef()
	} else {
		ref = fn.Ref()
	}

	return pairs(
		tmplFn,
		"in", stack.Indent(),
		"local", local,
		"ref", ref,
		"args", args.String(),
		"block", emitFnBlock(fn, varArgName),
	)
}

var (
	defaultInstanceName    = "_instance_"
	tmplDefaultConInstance = "local " + defaultInstanceName + " = setmetatable({}, self)"
	tmplVarArg             = `{in}local {ref} = \{ ... }`
	// Matches the Lua `self` keyword.
	// This is used to set a formatter hook while emitting an fn block. This
	// particular hook replaces all values passed into the formatter which match
	// the pattern with the default instance name `_instance_`.
	//
	// This ultimately allows us to create constructors using the same `this`
	// reference as regular methods, but
	selfRepl = regexp.MustCompile(`\bself\b`)
)

func emitFnBlock(fn *typeset.Fn, vararg string) string {
	f := formatter.NewFormatter()
	stack.Push()
	defer stack.Pop()

	// If the Fn is an overridden constructor, create the boilerplate
	// `_instance_`. Also set a formatter hook to replace all occurrences of
	// `this` in the constructor to the boilerplate `_instance_` object we create.
	if fn.Constructor {
		f.Newline().
			Str(stack.Indent()).
			Str(tmplDefaultConInstance)

		// Set the hook to replace `this` references (see comment above).
		unset := f.Hook(
			func(s string) string {
				return selfRepl.ReplaceAllLiteralString(s, defaultInstanceName)
			})

		// Defer the unset of the hook.
		defer unset()
	}

	// Vararg capture.
	if vararg != "" {
		f.Newline().
			Str(
				pairs(
					tmplVarArg,
					"ref", vararg,
					"in", stack.Indent(),
				))
	}

	// Block
	stack.Pop() // Scoping is inverted here since `emitBlock` scopes as well.
	f.Str(
		emitBlock(fn.Block),
	)
	stack.Push() // Scoping is inverted here since `emitBlock` scopes as well.

	// Defers
	for _, d := range stack.s[stack.i] {
		f.Newline().
			Str(stack.Indent()).
			Str(d)
	}

	return f.String()
}

var tmplFnA = `
function({args}) return {stmt} end
`

func emitAnonFn(fn *typeset.Fn) string {
	// Args
	args := formatter.NewFormatter()
	for i, arg := range fn.Args {
		// Handle standard args.
		args.Str(arg.Ref())

		if i < len(fn.Args)-1 {
			args.Str(", ")
		}
	}

	// Block
	v := emitValues(fn.Values)

	return formatter.NewFormatter().
		Str(
			pairs(
				tmplFnA,
				"args", args.String(),
				"stmt", v,
			)).String()
}

/*------------------------------------------------------------------------------
 * Statements
 *----------------------------------------------------------------------------*/

func emitStatement(stmt *typeset.Statement, isDefer bool) string {
	switch stmt.StmtType {
	// 'defer'
	case token.Defer:
		for d := range stmt.Defers() {
			stack.Add(emitStatement(d, true))
		}
		return ""

	// Call
	case token.Call:
		return emitCall(stmt.Call, false, isDefer)

	// Bind | Rebind
	case token.Bind, token.Rebind:
		return emitBind(stmt.Bind)

	// 'fn'
	case token.Fn:
		return emitFn(stmt.Fn)

	// 'for'
	case token.For:
		return emitFor(stmt.For)

	// 'if'
	case token.If:
		return emitConditional(stmt.If)

	// Reference
	case token.Ref:
		return stmt.Ref()

	// 'return'
	case token.Ret:
		f := formatter.NewFormatter()
		// When we hit a `return`, unwind and emit any queued `defer`s.
		for d := range stack.Unwind() {
			// Don't add a newline to the first member.
			if d.current > 0 {
				f.Newline()
			}
			f.Str(stack.Indent()).
				Str(d.stmt)
		}

		// If we actually had defers to unwind, we need a newline between the final
		// deferred statement and the return statement.
		if stack.defersTotal > 0 {
			f.Newline()
		}

		f.Str(stack.Indent()).
			Str("return")

		// A minute detail, but by checking the length here we avoid emitting bare
		// returns with a trailing whitespace.
		if len(stmt.Values) > 0 {
			f.Space()
		}

		f.Str(
			emitValues(stmt.Values),
		)

		return f.String()

	default:
		sklog.UnexpectedType(
			"eStatement()",
			stmt.StmtType,
		)
		return "" // Unreachable
	}
}

/*------------------------------------------------------------------------------
 * Values
 *----------------------------------------------------------------------------*/

func emitValues(values []*typeset.Value) string {
	f := formatter.NewFormatter()

	for _, v := range values {
		f.Str(
			emitValue(v),
		)
	}

	return f.String()
}

func emitValue(value *typeset.Value) string {
	switch value.ValueType {
	// List literal.
	case token.ListL:
		return emitList(value.List)

	// '[]'
	case token.List:
		return "{}"

	// Anonymous fn
	case token.Fn:
		return emitAnonFn(value.Fn)

	// Value group
	case token.ValueGroup:
		return "(" + emitValues(value.Group) + ")"

	// Call
	case token.Call:
		return emitCall(value.Call, true, false)

	// Reference
	case token.Ref:
		return lua.Translate(value.Ref())

	// StrL
	case token.StrL:
		return "'" + value.StrL + "'"

	// IntL
	case token.IntL:
		return value.IntL

	// BoolL
	case token.BoolL:
		return value.BoolL

	// Nil
	case token.Nil:
		return value.Nil

	// '!'
	case token.Not:
		return "not "

	// Operators
	case token.MathOperator, token.ComparisonOperator, token.ConcatOperator, token.LogicOperator:
		return " " + lua.Translate(value.Op) + " "

	default:
		sklog.UnexpectedType("emit Value", value.ValueType)
		return ""
	}
}

func emitList(list []*typeset.Value) string {
	f := formatter.NewFormatter()

	f.Str("{")
	stack.Push()
	for i, listVal := range list {
		f.Newline().
			Str(stack.Indent()).
			Str(
				emitValue(listVal),
			)

		// Comma delimit all but the last list member.
		if i < len(list)-1 {
			f.Str(",")
		}
	}
	stack.Pop()

	return f.Newline().Str(stack.Indent()).Str("}").String()
}

/*------------------------------------------------------------------------------
 * For
 *----------------------------------------------------------------------------*/

var tmplForI = `
{in}for {iterators} = {iterables} do{block}
{in}{end}
`

var tmplForIn = `
{in}for {iterators} in pairs({iterables}) do{block}
{in}{end}
`

func emitFor(nfor *typeset.For) string {
	// Identify the template based on whether it's a `for k, v in` or
	// `for i = n, n` format loop.
	var tmpl string
	if len(nfor.Iterables) == 1 {
		tmpl = tmplForIn
	} else {
		tmpl = tmplForI
	}

	// Prepare the iterator(s).
	iterators := forIterators(nfor)

	// Prepare the iterable(s).
	iterables := forIterables(nfor)

	// Emit all contained statements.
	block := emitBlock(nfor.Block)

	return pairs(
		tmpl,
		"in", stack.Indent(),
		"iterators", iterators,
		"iterables", iterables,
		"block", block,
	)
}

func forIterators(iterators *typeset.For) string {
	f := formatter.NewFormatter()

	// Stringify the iterator(s).
	for i, v := range iterators.Iterators {
		f.Str(v.Ref())
		if i < len(iterators.Iterators)-1 {
			f.Str(", ")
		}
	}

	return f.String()
}

func forIterables(iterables *typeset.For) string {
	f := formatter.NewFormatter()

	// Stringify the iterable(s).
	for i, v := range iterables.Iterables {
		if v.Value != "" {
			f.Str(v.Value)
		} else {
			f.Str(v.Ref())
		}

		if i < len(iterables.Iterables)-1 {
			f.Str(", ")
		}
	}

	return f.String()
}

/*------------------------------------------------------------------------------
 * If
 *----------------------------------------------------------------------------*/

func emitConditional(cond *typeset.If) string {
	f := formatter.NewFormatter()

	// If
	if v := emitIf(cond); v != "" {
		f.Str(v)
	}

	// Elifs
	if v := emitElifs(cond.Elifs); v != "" {
		f.Str(v)
	}

	// Else
	if v := emitElse(cond.Else); v != "" {
		f.Str(v)
	}

	// Close the conditional.
	f.Newline().Str(stack.Indent()).Str("end")

	return f.String()
}

var tmplIf = `
{in}if {conds} then{block}`

func emitIf(nif *typeset.If) string {
	// Conditions
	conds := emitConditions(nif.Conditions)

	// Block
	block := emitBlock(nif.Block)

	return pairs(
		tmplIf,
		"in", stack.Indent(),
		"conds", conds,
		"block", block,
	)
}

var tmplElif = `
{in}elseif {conds} then{block}
`

func emitElifs(elifs []*typeset.Elif) string {
	if len(elifs) == 0 {
		return ""
	}

	f := formatter.NewFormatter()
	for _, s := range elifs {
		// Add the elif.
		f.Newline().Str(
			pairs(
				tmplElif,
				"in", stack.Indent(),
				"conds", emitConditions(s.Conditions),
				"block", emitBlock(s.Block),
			))
	}

	return f.String()
}

var tmplElse = `
{in}else{block}
`

func emitElse(nelse *typeset.Else) string {
	if nelse == nil {
		return ""
	}

	f := formatter.NewFormatter()
	f.Newline().Str(
		pairs(
			tmplElse,
			"in", stack.Indent(),
			"block", emitBlock(nelse.Block),
		))

	return f.String()
}

func emitConditions(conds []*typeset.Value) string {
	f := formatter.NewFormatter()

	for _, c := range conds {
		f.Str(emitValue(c))
	}

	return f.String()
}

/*------------------------------------------------------------------------------
 * Blocks
 *----------------------------------------------------------------------------*/

func emitBlock(block []*typeset.Statement) string {
	f := formatter.NewFormatter()
	stack.Push()
	defer stack.Pop()

	// Emit all contained statements.
	for _, v := range block {
		if stmt := emitStatement(v, false); stmt != "" {
			f.Newline().Str(emitStatement(v, false))
		}
	}

	// Write any defers at the close of the block.
	for _, d := range stack.s[stack.i] {
		f.Newline().Str(stack.Indent()).Str(d)
	}

	return f.String()
}

/*------------------------------------------------------------------------------
 * Binds
 *----------------------------------------------------------------------------*/

var tmplDecl = `
{in}{local}{ref}
`
var tmplBind = `
{in}{local}{ref} = {Value}
`

func emitBind(bind *typeset.Bind) string {
	// Indentation.
	indent := stack.Indent()

	// Handle access scoping.
	// Binding is not declared pub.
	var local string
	if !bind.Pub() &&
		// Binding is not a rebind of an existing Value.
		!bind.Rebind {
		local = "local "
	}

	bindsF := formatter.NewFormatter()
	for i, boundRef := range bind.Binds {
		bindsF.Str(boundRef.Ref())

		if i < len(bind.Binds)-1 {
			bindsF.Str(", ")
		}
	}

	// If we have no values, return here (decl only).
	if len(bind.Values) == 0 {
		return pairs(
			tmplDecl,
			"in", indent,
			"local", local,
			"ref", bindsF.String(),
		)
	}

	// Prepare the values.
	values := formatter.NewFormatter()
	for _, v := range bind.Values {
		values.Str(emitValue(v))
	}

	// String it all together.
	return formatter.NewFormatter().Str(
		pairs(
			tmplBind,
			"in", indent,
			"local", local,
			"ref", bindsF.String(),
			"Value", values.String(),
		)).String()
}

/*------------------------------------------------------------------------------
 * Calls
 *----------------------------------------------------------------------------*/

var tmplCall = `
{in}{ref}({args}){newline}
`

func emitCall(call *typeset.Call, value bool, isDefer bool) string {
	if call.Ref() == "" && len(call.Args) == 0 {
		return ""
	}

	// Indentation
	var indent string
	var newline string
	if !value && !isDefer {
		indent = stack.Indent()
		newline = "\n"
	}

	// Reference
	var ref string
	if call.RefsLen() > 1 {
		ref = call.MethodRef()
	} else {
		ref = call.Ref()
	}
	ref = lua.Translate(ref)

	// Args
	args := formatter.NewFormatter()
	for i, arg := range call.Args {
		for _, v := range arg.Values {
			if arg.Spread {
				args.Str("...")
			} else {
				args.Str(emitValue(v))
			}
		}
		if i < len(call.Args)-1 {
			args.Str(", ")
		}
	}

	return pairs(
		tmplCall,
		"in", indent,
		"ref", ref,
		"args", args.String(),
		"newline", newline,
	)
}
