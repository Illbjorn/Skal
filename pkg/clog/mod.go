package clog

import "os"

var (
	exit = os.Exit
)

const (
	bSpace   = ' '
	bNewline = '\n'
	bQuote   = '\''
	bEq      = '='
)

var (
	bArrSpace   = []byte{bSpace}
	bArrNewline = []byte{bNewline}
	bArrQuote   = []byte{bQuote}
	bArrEq      = []byte{bEq}
)
