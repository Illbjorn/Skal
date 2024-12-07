package time

import (
	"net/http"
	"time"

	lua "github.com/yuin/gopher-lua"
)

const moduleName = "time"

func Load(l *lua.LState) {
	// Fns
	t := l.SetFuncs(l.NewTable(), fns())

	// Vars
	for k, v := range vars() {
		l.SetField(t, k, v)
	}

	// Set the global var.
	l.SetGlobal(moduleName, t)
}

func fns() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"now": now,
	}
}

func vars() map[string]lua.LValue {
	return map[string]lua.LValue{}
}

var (
	client = http.Client{Timeout: 2000 * time.Millisecond}
)
