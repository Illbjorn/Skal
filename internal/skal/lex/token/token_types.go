package token

type Type uint8

func (t Type) String() string {
	switch t {
	// Keywords
	case This:
		return "this"
	case Pub:
		return "pub"
	case New:
		return "new"
	case Ret:
		return "return"
	case For:
		return "for"
	case In:
		return "in"
	case Import:
		return "import"
	case Defer:
		return "defer"
	case Fn:
		return "fn"
	case Struct:
		return "struct"
	case Enum:
		return "enum"
	case Let:
		return "let"

	// Primitive Types
	case Int:
		return "int"
	case Bool:
		return "bool"
	case Str:
		return "str"
	case Undefined:
		return "undefined"

	// Extern
	case Extern:
		return "extern"
	case As:
		return "as"

	// Conditional Keywords
	case If:
		return "if"
	case Elif:
		return "elif"
	case Else:
		return "else"

	// Literals
	case True:
		return "true"
	case False:
		return "false"
	case Nil:
		return "nil"
	case List:
		return "[]"

	// Comparison Operators
	case EQEQ:
		return "=="
	case GE:
		return ">="
	case GT:
		return ">"
	case LT:
		return "<"
	case LE:
		return "<="
	case NE:
		return "!="

	// Logic Operators
	case Not:
		return "!"
	case And:
		return "&&"
	case Or:
		return "||"

	// String Concatenation
	case Concat:
		return ".."

	// Math Operators
	case Plus:
		return "+"
	case Minus:
		return "-"
	case Mult:
		return "*"
	case Div:
		return "/"

	// Spread Operator
	case Spread:
		return "..."

	// Ungrouped
	case Colon:
		return ":"
	case Space:
		return " "
	case LF:
		return "\n"
	case Comment:
		return "#"
	case Dot:
		return "."
	case EQ:
		return "="
	case SemiColon:
		return ";"
	case BraceOpen:
		return "{"
	case BraceClose:
		return "}"
	case BrackOpen:
		return "["
	case BrackClose:
		return "]"
	case ParenOpen:
		return "("
	case ParenClose:
		return ")"
	case Comma:
		return ","
	case Arrow:
		return "->"

	// Type System
	case TypeHint:
		return "type hint"

	// Generic Collection
	case Invalid:
		return "invalid"
	case EOF:
		return "<EOF>"

	// Extern
	case External:
		return "external"

	// Literals
	case IntL:
		return "int literal"
	case StrL:
		return "str literal"
	case BoolL:
		return "bool literal"
	case ListL:
		return "list literal"

	// Operator Types
	case ComparisonOperator:
		return "comparison operator"
	case LogicOperator:
		return "logic operator"
	case ConcatOperator:
		return "concat operator"
	case MathOperator:
		return "math operator"

	// Conditionals
	case Elifs:
		return "elifs"

	// Naming | References
	case ID:
		return "identifier"
	case Ref:
		return "reference"
	case Index:
		return "index"

	// Values
	case Value:
		return "Value"
	case ValueGroup:
		return "Value group"

	// Statements
	case Statement:
		return "statement"
	case Block:
		return "block"

	// Conditions
	case Conditions:
		return "conditions"
	case ConditionsGroup:
		return "conditions group"

	// Calls
	case Call:
		return "call"
	case CallArg:
		return "call arg"

	// Structs
	case StructField:
		return "struct field"
	case StructFields:
		return "struct fields"
	case StructMethod:
		return "struct method"

	// Enums
	case EnumMember:
		return "enum member"
	case EnumMembers:
		return "enum members"

	// Functions
	case AFn:
		return "anonymous fn"
	case FnArg:
		return "fn arg"

	// Binds
	case Bind:
		return "bind"
	case Rebind:
		return "rebind"

	// For
	case ForType:
		return "for type"
	case ForIterable:
		return "for iterable"
	case ForIterator:
		return "for iterator"

	default:
		return ""
	}
}

const (
	//////////////////////////////// TOKEN VALUES ////////////////////////////////
	//                                                                          //
	//                        These are parsed literally.                       //
	//                                                                          //
	//////////////////////////////////////////////////////////////////////////////
	//
	// Keywords
	This   Type = iota // Struct instance self-reference.
	Pub                // Public modifier.
	New                // Struct constructor.
	Ret                // Fn return.
	For                // For loop.
	In                 // For 'in'.
	Import             // Module import statement.
	Defer              // Defer statement.
	Fn                 // Function definition.
	Struct             // Struct definition.
	Enum               // Enum definition.
	Let                // Let binding.

	// Primitive Types
	Int
	Bool
	Str
	Undefined

	// Extern
	Extern
	As

	// Conditional Keywords
	If
	Elif
	Else

	// Literals
	True
	False
	Nil
	List // List Literal

	// Comparison Operators
	EQEQ
	GE
	GT
	LT
	LE
	NE

	// Logic Operators
	Not
	And
	Or

	// String Concatenation
	Concat

	// Math Operators
	Plus
	Minus
	Mult
	Div

	// Spread Operator
	Spread

	// Ungrouped
	Colon      // Type Hint
	Space      // Whitespace
	LF         // Newline
	Comment    // Src Comment
	Dot        // Accessor
	EQ         // Assignment
	SemiColon  // Src Terminator
	BraceOpen  // Scope Begin
	BraceClose // Scope End
	BrackOpen  // Index Begin
	BrackClose // Index End
	ParenOpen  // Group Begin
	ParenClose // Group End
	Comma      // Punctuation, separator
	Arrow      // Anonymous function

	//////////////////////////// CATEGORIZATION VALUES ///////////////////////////
	//                                                                          //
	//These only serve to indicate a broader grouping associated with the token.//
	//                                                                          //
	//////////////////////////////////////////////////////////////////////////////

	// Type System
	TypeHint

	// Generic Collection
	Invalid // Returned by the Tokenizer to signal no lookahead or behind.
	EOF     // Returned by the Tokenizer to signal end of lexing.

	// Extern
	External // An individual foreign reference.

	// Literals
	IntL  // Integer literal.
	StrL  // String literal.
	BoolL // Boolean true/false.
	ListL // List literal.

	// Operator Types
	ComparisonOperator
	LogicOperator
	ConcatOperator
	MathOperator

	// Conditionals
	Elifs // Elifs are categorized under an umbrella for AST construction.

	// Naming | References
	ID    // Example: my_struct, my_fn
	Ref   // Examples: my_struct.MyField, this.MyField
	Index // Example: my_struct.MyArrField[1]

	// Values
	Value      // Examples: true, 12 + 3, false
	ValueGroup // Example: 12 + 3

	// Statements
	Statement // For, if, fn, etc.
	Block     // Group of `tStatement`.

	// Conditions
	Conditions      // Really just a `tValues`. Separate for AST purposes.
	ConditionsGroup // Really just a `tValueGroup`. Separate for AST purposes.

	// Calls
	Call    // Fn or method call.
	CallArg // Fn or method call argument.

	// Structs
	StructField
	StructFields
	StructMethod

	// Enums
	EnumMember
	EnumMembers

	// Functions
	AFn
	FnArg

	// Binds
	Bind   // Let binding.
	Rebind // Reassignment of a previously bound variable.

	// For
	ForType     // 'in' | '='
	ForIterable // ex: The `Value` in: for k, v in Value {
	ForIterator // ex: The `i` in: for i = 1, 10 {
)
