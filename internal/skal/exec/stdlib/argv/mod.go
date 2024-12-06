package argv

import (
	"log/slog"
	"strconv"
)

var (
	itoa = strconv.Itoa
)

func dbgLogger(loc string) *slog.Logger {
	// stack := callStackS(3)
	l := slog.With("loc", loc)
	// for i := range stack {
	// 	l = l.With("call_stack_"+itoa(i), stack[i])
	// }

	return l
}
