package clog

var flags clogFlags

type clogFlags uint8

const (
	withShortTime clogFlags = 1 << iota
	withLongTime
	withFile
	withLevel
)

func WithFile() {
	flags = flags | withFile
}

func WithLevel() {
	flags = flags | withLevel
}
