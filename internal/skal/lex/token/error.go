package token

import (
	"github.com/illbjorn/skal/internal/skal/sklog"
	"github.com/illbjorn/skal/pkg/fstr"
)

func assertError(src string, tk Token, msg string) {
	file := tk.File()
	line := tk.LineStart()
	col1 := tk.ColumnStart()
	col2 := tk.ColumnEnd()
	sklog.NewCompilerEvent(sklog.MsgTypeParseError, sklog.LevelFatal).
		WithCallStack(3).
		WithSourceHint(src, file, line, col1, col2).
		Str(msg).
		Send()
}

func assertErrorF(src string, tk Token, msg string, pairs ...string) {
	assertError(src, tk, fstr.PairsStrip(msg, pairs...))
}
