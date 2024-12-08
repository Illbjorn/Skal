package clog

type level uint8

const (
	lDebug level = iota
	lInfo
	lWarn
	lError
	lFatal
)

func (l level) String() string {
	switch l {
	case lDebug:
		return "DBG"
	case lInfo:
		return "INF"
	case lWarn:
		return "WRN"
	case lError:
		return "ERR"
	case lFatal:
		return "FTL"
	default:
		return ""
	}
}
