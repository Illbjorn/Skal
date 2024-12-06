package stdlib

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/http"
	lua "github.com/yuin/gopher-lua"
)

type loader func(l *lua.LState)

var libLoaders = []loader{
	http.Load,
	conv.Load,
}

func Load(l *lua.LState) {
	for _, loader := range libLoaders {
		loader(l)
	}
}
