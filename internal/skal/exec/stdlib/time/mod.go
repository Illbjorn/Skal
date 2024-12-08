package time

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/module"
	lua "github.com/yuin/gopher-lua"
)

func Module() module.Module {
	return module.Module{
		Name: "time",
		ModuleFns: map[string]lua.LGFunction{
			"now": now,
		},
	}
}
