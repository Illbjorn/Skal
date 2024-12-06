package conv

import (
	"encoding/json"
	"log/slog"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/argv"
	lua "github.com/yuin/gopher-lua"
)

func toJSON(l *lua.LState) int {
	var (
		v = l.Get(-1)
	)

	if !argv.IsTable(v) {
		return 0
	}

	m := make(map[string]any)

	v.(*lua.LTable).ForEach(
		func(l1, l2 lua.LValue) {
			m[l1.String()] = l2
		},
	)

	content, err := json.Marshal(m)
	if err != nil {
		slog.Debug(
			"Failed to marshal JSON content.",
			"error", err,
		)

		return 0
	}

	l.Push(lua.LString(string(content)))

	return 1
}
