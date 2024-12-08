package clog

import (
	"os"
)

func init() {
	// WithFile()
	// WithLevel()
}

func log(l level, v string, pairs ...any) {

	fmtPrefix(l)

	os.Stdout.Write([]byte(v))

	if len(pairs) > 0 {
		os.Stdout.Write(bArrSpace)
		fmtPairs(pairs...)
	}

	os.Stdout.Write(bArrNewline)
}

func Debug(v string, pairs ...any) {
	log(lDebug, v, pairs...)
}

func Info(v string, pairs ...any) {
	log(lInfo, v, pairs...)
}

func Warn(v string, pairs ...any) {
	log(lWarn, v, pairs...)
}

func Error(v string, pairs ...any) {
	log(lError, v, pairs...)
}

func Fatal(v string, pairs ...any) {
	log(lFatal, v, pairs...)
	exit(1)
}
