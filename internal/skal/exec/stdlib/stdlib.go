package stdlib

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/fmt"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/http"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/time"
	lua "github.com/yuin/gopher-lua"
)

type loader func(l *lua.LState)

var libLoaders = []loader{
	http.Load,
	conv.Load,
	fmt.Load,
	time.Load,
}

func Load(l *lua.LState) {
	for _, loader := range libLoaders {
		loader(l)
	}
}
