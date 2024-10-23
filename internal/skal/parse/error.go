package parse

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func parseError(msg string, tk token.Token, fatal bool) {
	level := sklog.LevelError
	if fatal {
		level = sklog.LevelFatal
	}

	src := tk.SrcLine()
	file := tk.File()
	line := tk.LineStart()
	col1 := tk.ColumnStart()
	col2 := tk.ColumnEnd()

	sklog.
		NewCompilerEvent(sklog.MsgTypeParseError, level).
		WithCallStack(3).
		WithSourceHint(src, file, line, col1, col2).
		Str(msg).
		Send()
}
