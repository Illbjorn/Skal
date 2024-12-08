package fmt

import (
	"fmt"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/argv"
	"github.com/illbjorn/skal/pkg/clog"
	lua "github.com/yuin/gopher-lua"
)

func sprintf(l *lua.LState) int {
	if l.GetTop() < 3 {
		clog.Debug(
			"received insufficient args in call to sprintf",
		)

		return 0
	}

	var (
		tmpl   string
		values []any
	)

	// Collect args.
	for arg := range argv.Get(l) {
		if tmpl == "" {
			tmpl = arg.String()
			continue
		}

		values = append(values, arg)
	}

	if len(values) == 0 {
		clog.Debug(
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
	for arg := range argv.Get(l) {
		if tmpl == "" {
			tmpl = arg.String()
			continue
		}

		values = append(values, arg)
	}

	if len(values) == 0 {
		clog.Debug(
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
