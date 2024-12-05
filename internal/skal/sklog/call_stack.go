package sklog

import (
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/illbjorn/fstr"
	"github.com/illbjorn/skal/pkg/pprint"
)

// Produces a formatted string representation of the `i`th caller down the
// call stack.
func caller(i int) string {
	// Get the caller info.
	_, f, l, _ := runtime.Caller(i)

	// Cut down the file path to just the base file name.
	fn := filepath.Base(f)

	// Format and return.
	return fstr.Pairs(
		"{file}:{line}",
		"file", fn,
		"line", strconv.Itoa(l),
	)
}

// Produces a simple string slice representation of the current call stack.
//
// This function will iterate call stack depth from 1-10 looking for the first
// which resides outside of this package. This will serve as the starting index
// to build the call stack FROM to provided level `depth`.
func callStack(depth int) []string {
	// Walk down the stack until we find a file that is NOT in our current
	// package, then use that index+1 as the starting frame to generate our stack
	// trace.
	var from int
	for i := 0; i < 10; i++ {
		_, f, _, _ := runtime.Caller(i)
		base := filepath.Base(f)
		if !slices.Contains(
			[]string{
				"call_stack.go",
				"compiler_event.go",
				"error.go",
			}, base) {
			from = i + 1
			break
		}
	}

	var out []string
	// Assemble the stack trace.
	for i := from; i < from+depth; i++ {
		out = append(out, caller(i))
	}

	return out
}

// Wraps `callStack()` call output into a *pprint.Msg for use in log messages.
func callStackF(msg *pprint.Msg, depth int) *pprint.Msg {
	frames := callStack(depth)
	for i, v := range frames {
		// Add the 'AT' block in white.
		msg.Add("\n" + strings.Repeat(" ", i+1))
		msg.Yellow("AT: ")
		// Add the call stack value in yellow.
		msg.White(v)
	}
	return msg
}
