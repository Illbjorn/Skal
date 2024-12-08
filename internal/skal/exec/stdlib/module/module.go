package module

import (
	"github.com/illbjorn/skal/pkg/clog"
	lua "github.com/yuin/gopher-lua"
)

type Module struct {
	GlobalFns  map[string]lua.LGFunction
	ModuleFns  map[string]lua.LGFunction
	GlobalVars map[string]lua.LValue
	ModuleVars map[string]lua.LValue
	Name       string
}

func (m Module) Load(l *lua.LState) {
	var (
		pairs = []any{"module", m.Name}
	)

	clog.Debug("<-- beginning module load", pairs...)
	defer clog.Debug("module load complete --!>", pairs...)

	// Global Fns
	for k, v := range m.GlobalFns {
		clog.Debug("loading global fn", "fn", k)
		l.SetGlobal(k, l.NewFunction(v))
	}

	// Global Vars
	for k, v := range m.GlobalVars {
		clog.Debug("loading global var", "var", k)
		l.SetGlobal(k, v)
	}

	// Module Items
	if len(m.ModuleFns)+len(m.ModuleVars) > 0 {
		t := new(lua.LTable)

		// Module Fns
		for k, v := range m.ModuleFns {
			clog.Debug("loading module fn", "fn", k)
			t.RawSetString(k, l.NewFunction(v))
		}

		// Module Vars
		for k, v := range m.ModuleVars {
			clog.Debug("loading module var", "var", k)
			t.RawSetString(k, v)
		}

		l.SetGlobal(m.Name, t)
	}
}
