package sklog

import (
	"github.com/illbjorn/skal/pkg/fstr"
)

// Logs a generic fatal compiler error.
func CFatal(msg string) {
	NewCompilerEvent(MsgTypeCompilerError, LevelFatal).
		WithCallStack(3).
		Str(msg).
		Send()
}

// Adds a `fstr.PairsStrip()` message format before sending to `CFatal()`.
func CFatalF(msg string, pairs ...string) {
	CFatal(
		fstr.Pairs(msg, pairs...),
	)
}

// Produces a warning indicating a particular code branch or feature is not yet
// implemented.
//
//goland:noinspection GoUnusedExportedFunction
func Todo(loc string) {
	NewCompilerEvent(MsgTypeTodo, LevelWarn).
		WithCallStack(1).
		Str(loc).
		Send()
}

// Used for producing a generic terminating error where we expected a
// particular: token, node or emit statement type.
func UnexpectedType(loc, found string) {
	CFatalF(
		"Found unexpected {loc} type: {found}.",
		"loc", loc,
		"found", found,
	)
}
