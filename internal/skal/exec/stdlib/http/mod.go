package http

import (
	"net/http"
	"time"

	lua "github.com/yuin/gopher-lua"
)

func Load(l *lua.LState) {
	// Fns
	t := l.SetFuncs(l.NewTable(), fns())

	// Vars
	for k, v := range vars() {
		l.SetField(t, k, v)
	}

	// Set the global var.
	l.SetGlobal("http", t)
}

func fns() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"get":  get,
		"post": post,
	}
}

func vars() map[string]lua.LValue {
	return map[string]lua.LValue{}
}

var (
	client = http.Client{Timeout: 2000 * time.Millisecond}
)
