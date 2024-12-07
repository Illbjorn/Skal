//go:build linux || darwin

package time

import (
	"time"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	lua "github.com/yuin/gopher-lua"
)

type Time struct {
	Day, Hour, Minute, Second, Nanosecond int
	Millisecond                           int64
}

func now(l *lua.LState) int {
	n := time.Now()

	t := Time{
		Day:         n.Day(),
		Hour:        n.Hour(),
		Minute:      n.Minute(),
		Second:      n.Second(),
		Millisecond: n.UnixMilli(),
		Nanosecond:  n.Nanosecond(),
	}

	l.Push(conv.StructToLTable(t, l))

	return 1
}
