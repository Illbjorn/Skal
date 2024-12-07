package conv

import (
	"encoding/json"
	"fmt"
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

	l.Push(lua.LString(ToJSON(v.(*lua.LTable))))

	return 1
}

func ToJSON(t *lua.LTable) string {
	var (
		m = make(map[string]any)
	)

	t.ForEach(
		func(l1, l2 lua.LValue) {
			m[l1.String()] = l2
		},
	)

	var (
		content []byte
		err     error
	)

	if content, err = json.Marshal(m); err != nil {
		slog.Debug(
			"Failed to marshal JSON content.",
			"error", err,
		)

		return ""
	}

	return string(content)
}

func FromJSON(l *lua.LState) int {
	var (
		v = l.Get(-1)
	)

	if !argv.IsString(v) {
		return 0
	}

	t := l.NewTable()

	m := make(map[string]any)
	err := json.Unmarshal([]byte(v.String()), &m)
	if err != nil {
		slog.Debug(
			"failed to unmarshal JSON payload",
			"error", err,
		)

		return 0
	}

	t = buildJSON(t, l, m)

	l.Push(t)

	return 1
}

func buildJSON(t *lua.LTable, l *lua.LState, m map[string]any) *lua.LTable {
	for k, v := range m {
		switch v := v.(type) {
		case int:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case int8:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case int16:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case int32:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case int64:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case uint:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case uint8:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case uint16:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case uint32:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case uint64:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case float32:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case float64:
			l.SetField(t, k, lua.LNumber(float64(v)))

		case string:
			l.SetField(t, k, lua.LString(v))

		case bool:
			l.SetField(t, k, lua.LBool(v))

		default:
			slog.Error(
				"found unexpected JSON value type " + fmt.Sprintf("%T", v),
			)
		}
	}

	return t
}
