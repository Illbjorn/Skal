package argv

import (
	"bytes"
	"path/filepath"
	"runtime"
	"slices"
	"strconv"
	"strings"

	"github.com/illbjorn/fstr"
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

// Produces a simple string representation of the current call stack.
//
// This function will iterate call stack depth from 1-10 looking for the first
// which resides outside of the `log` prefixed Go files. The first level in the
// call stack which satisfies this condition will serve as the starting index
// to build the call stack FROM to provided count `depth`.
//
// # Example:
//
// Call stack:
//
//	 0: log_caller.go:44
//	-1: log_caller.go:80
//	-2: log.go:99
//	-3: validation.go:88
//	-4: get.go:111
//	-5: mod.go:123
//
// Call:
//
//	callStack(3)
//
// Result:
//
//	AT: validation.go:88
//	  AT: get.go:111
//	    AT: mod.go:123
func callStack(depth int) string {
	// Walk down the stack until we find a file that is NOT in our current
	// package, then use that index+1 as the starting frame to generate our stack
	// trace.
	var from int
	for i := 0; i < 10; i++ {
		_, f, _, _ := runtime.Caller(i)
		base := filepath.Base(f)
		if !slices.Contains(
			[]string{
				"log_caller.go",
				"log.go",
			}, base) {
			from = i + 1
			break
		}
	}

	ret := bytes.NewBuffer(nil)
	// Assemble the stack trace.
	for i := from; i < from+depth; i++ {
		ret.WriteString("\n" + strings.Repeat(" ", i+1))
		ret.WriteString("AT: " + caller(i))
	}

	return ret.String()
}

func callStackS(depth int) []string {
	// Walk down the stack until we find a file that is NOT in our current
	// package, then use that index+1 as the starting frame to generate our stack
	// trace.
	var from int
	for i := 0; i < 10; i++ {
		_, f, _, _ := runtime.Caller(i)
		base := filepath.Base(f)
		if !slices.Contains(
			[]string{
				"log_caller.go",
				"log.go",
			}, base) {
			from = i + 1
			break
		}
	}

	ret := make([]string, depth)
	// Assemble the stack trace.
	for i := from; i < from+depth; i++ {
		ret[from+depth-i-1] = caller(i)
	}

	return ret
}
