package validate

import (
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func verr(msg string, token token.Token) {
	event := sklog.NewCompilerEvent(sklog.MsgTypeValidationError, sklog.LevelError).
		WithCallStack(3).
		Str(msg)

	if token != nil {
		file := token.File()
		line := token.LineStart()
		col1 := token.ColumnStart()
		col2 := token.ColumnEnd()
		event = event.WithSourceHint(token.SrcLine(), file, line, col1, col2)
	}

	event.Send()
}
