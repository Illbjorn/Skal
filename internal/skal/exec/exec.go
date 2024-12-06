package exec

import (
	"github.com/illbjorn/skal/internal/skal/exec/stdlib"
	"github.com/illbjorn/skal/internal/skal/sklog"
	lua "github.com/yuin/gopher-lua"
)

func Exec(f string) {
	// Init Lua VM.
	l := lua.NewState()
	defer l.Close()

	// Load Stdlib.
	stdlib.Load(l)

	// Execute script.
	if err := l.DoString(f); err != nil {
		sklog.CFatal(err.Error())
	}
}
