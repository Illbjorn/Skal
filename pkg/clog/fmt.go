package clog

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

func fmtPairs(pairs ...any) {
	if len(pairs)%2 != 0 {
		return
	}

	if len(pairs) == 0 {
		return
	}

	for i := 0; i < len(pairs); i += 2 {
		var (
			anyK = pairs[i]
			anyV = pairs[i+1]
		)

		if v, ok := pairs[i].(string); ok {
			os.Stdout.Write([]byte(v))
		} else {
			fmtAnyToString(anyK)
		}

		os.Stdout.Write(bArrEq)

		if v, ok := pairs[i+1].(string); ok {
			os.Stdout.Write([]byte(v))
		} else {
			fmtAnyToString(anyV)
		}

		if i < len(pairs)-1 {
			os.Stdout.Write(bArrSpace)
		}
	}
}

func fmtPrefix(l level) {
	if flags == 0 {
		return
	}

	// if WithShortTime&flags == WithShortTime {
	// 	ret = append(ret, withShortTime())
	// }

	// if WithLongTime&flags == WithLongTime {
	// 	ret = append(ret, withLongTime())
	// }

	if withFile&flags == withFile {
		fmtWithFile()
	}

	if withLevel&flags == withLevel {
		fmtWithLevel(l)
	}
}

var pkgFiles = []string{
	"clog.go",
	"conv.go",
	"flags.go",
}

func fmtWithFile() {
	for i := 1; i < 5; i++ {
		var (
			_, f, l, _ = runtime.Caller(i)
		)
		f = filepath.Base(f)

		if !slices.Contains(pkgFiles, f) {
			os.Stdout.Write([]byte(f))
			os.Stdout.Write([]byte{':'})
			fmtBits(uint64(l), false)
			return
		}
	}
}

func fmtWithLevel(l level) {
	os.Stdout.Write([]byte(l.String()))
}
