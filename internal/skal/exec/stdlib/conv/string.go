package conv

import (
	"fmt"
	"log/slog"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func toString(l *lua.LState) int {
	var (
		v = l.Get(l.GetTop())
	)
	l.Pop(1)

	if v.Type() == lua.LTNil {
		return 0
	}

	l.Push(
		lua.LString(
			LVToString(v),
		),
	)

	return 1
}

func LVToString(value lua.LValue) string {
	switch v := value.(type) {
	case lua.LBool, lua.LString, lua.LNumber:
		return v.String()

	case *lua.LTable:
		return LTableToString(v, 0)

	default:
		slog.Debug(
			"found unexpected value type in ToString",
			"type", fmt.Sprintf("%T", v),
		)

		return ""
	}
}

func LTableToString(t *lua.LTable, depth int) string {
	var (
		pairs = make([]string, 0)
		ind   = strings.Repeat(" ", depth)
	)

	t.ForEach(
		func(l1, l2 lua.LValue) {
			var (
				k = l1.String()
				v string
			)

			switch l2 := l2.(type) {
			case lua.LBool:
				if l2 {
					v = "true"
				}
				v = "false"

			case lua.LString, lua.LNumber:
				v = l2.String()

			case *lua.LTable:
				pairs = append(
					pairs,
					ind+k+" = {\n"+LTableToString(l2, depth+1)+"\n"+ind+"}",
				)
				return

			default:
				slog.Debug(
					"found unexpected TableToString value type",
					"type", fmt.Sprintf("%T", v),
				)
				return
			}

			if strings.ContainsAny(v, " \n\t\f\v") {
				v = "'" + v + "'"
			}

			pairs = append(pairs, ind+k+" = "+v)
		},
	)

	return strings.Join(pairs, "\n")
}
