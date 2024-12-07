package fmt

import (
	"fmt"
	"log/slog"

	lua "github.com/yuin/gopher-lua"
)

func sprintf(l *lua.LState) int {
	var (
		tmpl   string
		values []any
	)

	// Collect args.
	for i := 2; i < l.GetTop(); i++ {
		var (
			v = l.Get(i)
		)

		fmt.Println("i:", i, "v:", v)

		if i == 2 {
			tmpl = v.String()
			continue
		}

		values = append(values, v)
	}
	l.Pop(l.GetTop() - 2)

	if len(values) == 0 {
		slog.Debug(
			"received zero-length values to sprintf",
		)
	}

	l.Push(
		lua.LString(
			fmt.Sprintf(
				tmpl,
				values...,
			),
		),
	)

	return 1
}

func printfln(l *lua.LState) int {
	var (
		tmpl   string
		values []any
	)

	// Collect args.
	for i := 2; i <= l.GetTop(); i++ {
		var (
			v = l.Get(i)
		)

		if v.Type() == lua.LTNil {
			continue
		}

		if i == 2 {
			tmpl = v.String()
			continue
		}

		values = append(values, v)
	}
	l.Pop(l.GetTop() - 2)

	if len(values) == 0 {
		slog.Debug(
			"received zero-length values to printfln",
		)
	}

	fmt.Println(
		fmt.Sprintf(
			tmpl, values...,
		),
	)

	return 0
}
