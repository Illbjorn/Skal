package lua

var StdlibFns = []string{
	"print",
	"insert",
	"random",
	"strmatch",
	"tonumber",
	"tostring",
	"date",
}

// Translates a provided Skal term to a Lua equivalent, if one exists.
func Translate(v string) string {
	switch v {
	// ----------------------------------------------------------------------------
	// Keywords
	case "this":
		return "self"

	// ----------------------------------------------------------------------------
	// Functions
	case "insert":
		return "table.insert"

	case "random":
		return "math.random"

	// ----------------------------------------------------------------------------
	// Operators
	case "!":
		return "not"

	case "!=":
		return "~="

	case "&&":
		return "and"

	case "||":
		return "or"

	// ----------------------------------------------------------------------------
	// N/A
	//
	default:
		return v
	}
}
