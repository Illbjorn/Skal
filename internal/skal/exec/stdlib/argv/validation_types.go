package argv

import (
	lua "github.com/yuin/gopher-lua"
)

// -----------------------------------------------------------------------------
// String
func IsString(v lua.LValue) bool {
	if v.Type() == lua.LTString {
		return true
	}

	l := dbgLogger("IsString")

	l.Debug(
		"Received invalid parameter value type.",
		"expected", "string",
		"received", v.Type().String(),
		"value", v.String(),
	)

	return false
}

// -----------------------------------------------------------------------------
// Table
func IsTable(v lua.LValue) bool {
	if v.Type() == lua.LTTable {
		return true
	}

	l := dbgLogger("IsTable")

	l.Debug(
		"Received invalid parameter value type.",
		"expected", "table",
		"received", v.Type().String(),
		"value", v.String(),
	)

	return false
}

// -----------------------------------------------------------------------------
// Nil
func IsNil(v lua.LValue) bool {
	if v.Type() == lua.LTNil {
		return true
	}

	l := dbgLogger("IsNil")

	l.Debug(
		"Received invalid parameter value type.",
		"expected", "nil",
		"received", v.Type().String(),
		"value", v.String(),
	)

	return false
}
