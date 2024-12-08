package stdlib

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/fmt"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/http"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/module"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/time"
	lua "github.com/yuin/gopher-lua"
)

type modInitFn func() module.Module

var modules = []modInitFn{
	http.Module,
	conv.Module,
	fmt.Module,
	time.Module,
}

func Load(l *lua.LState) {
	for _, fn := range modules {
		var (
			mod = fn()
		)

		mod.Load(l)
	}
}
