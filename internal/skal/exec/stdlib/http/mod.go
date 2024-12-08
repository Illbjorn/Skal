package http

import (
	"net/http"
	"time"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/module"
	lua "github.com/yuin/gopher-lua"
)

const moduleName = "http"

func Module() module.Module {
	return module.Module{
		Name: "http",
		ModuleFns: map[string]lua.LGFunction{
			"get":  get,
			"post": post,
		},
	}
}

var (
	client = http.Client{Timeout: 2000 * time.Millisecond}
)
