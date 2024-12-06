package argv

import (
	"net/http"
	"net/url"
	"slices"

	lua "github.com/yuin/gopher-lua"
)

// -----------------------------------------------------------------------------
// HTTP Request Method
var httpMethods = []string{
	http.MethodConnect,
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodOptions,
	http.MethodPatch,
	http.MethodPost,
	http.MethodPut,
	http.MethodTrace,
}

func IsHTTPMethod(v lua.LValue) bool {
	l := dbgLogger("IsHTTPMethod")

	if slices.Contains(httpMethods, v.String()) {
		return true
	}

	l.Debug(
		"Received invalid HTTP request method.",
		"value", v.String(),
	)

	return false
}

// -----------------------------------------------------------------------------
// Request URL
func IsURL(v lua.LValue) bool {
	l := dbgLogger("IsURL")

	if v.Type() != lua.LTString {
		l.Debug(
			"Received invalid parameter value type.",
			"expected", "string",
			"received", v.Type().String(),
			"value", v.String(),
		)

		return false
	}

	if len(v.String()) == 0 {
		l.Debug(
			"Received zero-length URL string.",
			"value", v.String(),
		)

		return false
	}

	var (
		err error
	)
	if _, err = url.Parse(v.String()); err != nil {
		l.Debug(
			"Failed to parse provided string value as a URL.",
			"error", err,
			"value", v.String(),
		)

		return false
	}

	return true
}
