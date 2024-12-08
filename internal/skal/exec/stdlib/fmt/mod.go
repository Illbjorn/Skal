package fmt

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/module"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "fmt"

func Module() module.Module {
	return module.Module{
		Name: "fmt",
		GlobalFns: map[string]lua.LGFunction{
			"printfln": printfln,
			"sprintf":  sprintf,
		},
	}
}
