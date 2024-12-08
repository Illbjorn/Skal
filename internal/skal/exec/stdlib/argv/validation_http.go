package argv

import (
	"net/http"
	"net/url"
	"slices"

	"github.com/illbjorn/skal/pkg/clog"
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
	if slices.Contains(httpMethods, v.String()) {
		return true
	}

	clog.Debug(
		"Received invalid HTTP request method.",
		"method", "IsHTTPMethod",
		"value", v.String(),
	)

	return false
}

// -----------------------------------------------------------------------------
// Request URL
func IsURL(v lua.LValue) bool {
	if v.Type() != lua.LTString {
		clog.Debug(
			"Received invalid parameter value type.",
			"method", "IsURL",
			"expected", "string",
			"received", v.Type().String(),
			"value", v.String(),
		)

		return false
	}

	if len(v.String()) == 0 {
		clog.Debug(
			"Received zero-length URL string.",
			"method", "IsURL",
			"value", v.String(),
		)

		return false
	}

	var (
		err error
	)
	if _, err = url.Parse(v.String()); err != nil {
		clog.Debug(
			"Failed to parse provided string value as a URL.",
			"method", "IsURL",
			"error", err,
			"value", v.String(),
		)

		return false
	}

	return true
}
