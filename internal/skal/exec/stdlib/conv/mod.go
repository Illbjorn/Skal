package conv

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/module"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "conv"

func Module() module.Module {
	return module.Module{
		Name: "conv",
		ModuleFns: map[string]lua.LGFunction{
			"to_json":   toJSON,
			"from_json": FromJSON,
			"to_string": toString,
		},
	}
}
