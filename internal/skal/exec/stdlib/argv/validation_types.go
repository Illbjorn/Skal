package argv

import (
	"github.com/illbjorn/skal/pkg/clog"
	lua "github.com/yuin/gopher-lua"
)

// -----------------------------------------------------------------------------
// String
func IsString(v lua.LValue) bool {
	if v.Type() == lua.LTString {
		return true
	}

	clog.Debug(
		"Received invalid parameter value type.",
		"method", "IsString",
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

	clog.Debug(
		"Received invalid parameter value type.",
		"method", "IsTable",
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

	clog.Debug(
		"Received invalid parameter value type.",
		"method", "IsNil",
		"expected", "nil",
		"received", v.Type().String(),
		"value", v.String(),
	)

	return false
}
