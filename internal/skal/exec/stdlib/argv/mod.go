package argv

import (
	"log/slog"
)

func dbgLogger(loc string) *slog.Logger {
	l := slog.With("loc", loc)

	return l
}
