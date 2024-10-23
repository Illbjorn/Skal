package token

//goland:noinspection GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst,GoUnusedConst
const (
	// ////////////////////////////// TOKEN VALUES ////////////////////////////////
	//                                                                          //
	//                        These are parsed literally.                       //
	//                                                                          //
	// ////////////////////////////////////////////////////////////////////////////
	// Keywords
	This   = "this"   // Struct instance self-reference.
	Pub    = "pub"    // Public modifier.
	New    = "new"    // Struct constructor.
	Ret    = "return" // Fn return.
	For    = "for"    // For loop.
	In     = "in"     // For 'in'.
	Import = "import" // Module import statement.
	Defer  = "defer"  // Defer statement.
	Fn     = "fn"     // Function definition.
	Struct = "struct" // Struct definition.
	Enum   = "enum"   // Enum definition.
	Let    = "let"    // Let binding.

	// Primitive Types
	Int       = "int"
	Bool      = "bool"
	Str       = "str"
	Undefined = "undefined"

	// Extern
	Extern = "extern"
	As     = "as"

	// Conditional Keywords
	If   = "if"
	Elif = "elif"
	Else = "else"

	// Literals
	True  = "true"
	False = "false"
	Nil   = "nil"
	List  = "[]" // List Literal

	// Comparison Operators
	EQEQ = "=="
	GE   = ">="
	GT   = ">"
	LT   = "<"
	LE   = "<="
	NE   = "!="

	// Logic Operators
	Not = "!"
	And = "&&"
	Or  = "||"

	// String Concatenation
	Concat = ".."

	// Math Operators
	Plus  = "+"
	Minus = "-"
	Mult  = "*"
	Div   = "/"

	// Spread Operator
	Spread = "..."

	// Ungrouped
	Colon      = ":"  // Type Hint
	Space      = " "  // Whitespace
	LF         = "\n" // Newline
	Comment    = "#"  // Src Comment
	Dot        = "."  // Accessor
	EQ         = "="  // Assignment
	SemiColon  = ";"  // Src Terminator
	BraceOpen  = "{"  // Scope Begin
	BraceClose = "}"  // Scope End
	BrackOpen  = "["  // Index Begin
	BrackClose = "]"  // Index End
	ParenOpen  = "("  // Group Begin
	ParenClose = ")"  // Group End
	Comma      = ","  // Punctuation, separator
	Arrow      = "->" // Anonymous function

	// ////////////////////////// CATEGORIZATION VALUES ///////////////////////////
	//                                                                          //
	// These only serve to indicate a broader grouping associated with the token.//
	//                                                                          //
	// ////////////////////////////////////////////////////////////////////////////

	// Type System
	TypeHint = "type hint"

	// Generic Collection
	Invalid = "invalid" // Returned by the Tokenizer to signal no lookahead or behind.
	EOF     = "<EOF>"   // Returned by the Tokenizer to signal end of lexing.

	// Extern
	External = "external" // An individual foreign reference.

	// Literals
	IntL  = "int literal"  // Integer literal.
	StrL  = "str literal"  // String literal.
	BoolL = "bool literal" // Boolean true/false.
	ListL = "list literal" // List literal.

	// Operator Types
	ComparisonOperator = "comparison operator"
	LogicOperator      = "logic operator"
	ConcatOperator     = "concat operator"
	MathOperator       = "math operator"

	// Conditionals
	Elifs = "elifs" // Elifs are categorized under an umbrella for AST construction.

	// Naming | References
	ID    = "identifier" // Example: my_struct, my_fn
	Ref   = "reference"  // Examples: my_struct.MyField, this.MyField
	Index = "index"      // Example: my_struct.MyArrField[1]

	// Values
	Value      = "Value"       // Examples: true, 12 + 3, false
	ValueGroup = "Value group" // Example: 12 + 3

	// Statements
	Statement = "statement" // For, if, fn, etc.
	Block     = "block"     // Group of `tStatement`.

	// Conditions
	Conditions      = "conditions"       // Really just a `tValues`. Separate for AST purposes.
	ConditionsGroup = "conditions group" // Really just a `tValueGroup`. Separate for AST purposes.

	// Calls
	Call    = "call"     // Fn or method call.
	CallArg = "call arg" // Fn or method call argument.

	// Structs
	StructField  = "struct field"
	StructFields = "struct fields"
	StructMethod = "struct method"

	// Enums
	EnumMember  = "enum member"
	EnumMembers = "enum members"

	// Functions
	AFn   = "anonymous fn"
	FnArg = "fn arg"

	// Binds
	Bind   = "bind"   // Let binding.
	Rebind = "rebind" // Reassignment of a previously bound variable.

	// For
	ForType     = "for type"     // 'in' | '='
	ForIterable = "for iterable" // ex: The `Value` in: for k, v in Value {
	ForIterator = "for iterator" // ex: The `i` in: for i = 1, 10 {
)
