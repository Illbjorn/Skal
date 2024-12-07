package fmt

import lua "github.com/yuin/gopher-lua"

const moduleName = "fmt"

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
		"printfln": printfln,
		"sprintf":  sprintf,
	}
}

func vars() map[string]lua.LValue {
	return map[string]lua.LValue{}
}
