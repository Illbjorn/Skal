package clog

import (
	"fmt"
	"os"
	"strconv"
)

var (
	fint   = strconv.FormatInt
	fuint  = strconv.FormatUint
	ffloat = strconv.FormatFloat
	fbool  = strconv.FormatBool
)

func fmtAnyToString(v any) {
	switch v := v.(type) {
	// String
	case string:
		os.Stdout.Write([]byte(v))

	// Signed Integers
	case int:
		fmtBits(uint64(v), v < 0)
	case int8:
		fmtBits(uint64(v), v < 0)
	case int16:
		fmtBits(uint64(v), v < 0)
	case int32:
		fmtBits(uint64(v), v < 0)
	case int64:
		fmtBits(uint64(v), v < 0)

	// Unsigned Integers
	case uint:
		fmtBits(uint64(v), v < 0)
	case uint8:
		fmtBits(uint64(v), v < 0)
	case uint16:
		fmtBits(uint64(v), v < 0)
	case uint32:
		fmtBits(uint64(v), v < 0)
	case uint64:
		fmtBits(uint64(v), v < 0)

	// Floats
	case float32:
		fmtBits(uint64(v), v < 0)
	case float64:
		fmtBits(uint64(v), v < 0)

	// Errors
	case error:
		os.Stdout.Write([]byte(v.Error()))

	// Stringer
	case fmt.Stringer:
		os.Stdout.Write([]byte(v.String()))

	default:
		return
	}
}
