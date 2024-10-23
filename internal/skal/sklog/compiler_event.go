package sklog

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/illbjorn/skal/pkg/fstr"
	"github.com/illbjorn/skal/pkg/pprint"
)

const (
	// CompilerEvent Levels
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelDebug = "DBG"
	LevelError = "ERR"
	LevelFatal = "FTL"
)

//goland:noinspection GoUnusedConst,GoUnusedConst
const (
	// CompilerEvent Message Types
	MsgTypeTodo            = "Todo"
	MsgTypeParseError      = "Parse Error"
	MsgTypeValidationError = "Validation Error"
	MsgTypeEmitError       = "Emit Error"
	MsgTypeTypesetError    = "Typeset Error"
	MsgTypeCompilerError   = "Compiler Error"
)

func NewCompilerEvent(mtype, level string) *CompilerEvent {
	return &CompilerEvent{
		mtype: mtype,
		level: level,
	}
}

type CompilerEvent struct {
	stack *pprint.Msg
	msg   *pprint.Msg
	hint  *pprint.Msg
	mtype string
	level string
}

func (m *CompilerEvent) WithCallStack(depth int) *CompilerEvent {
	m.stack = callStackF(pprint.New(), depth)
	return m
}

type SrcHintable interface {
	File() string
	LineStart() int
	ColumnStart() int
	ColumnEnd() int
}

func (m *CompilerEvent) WithSourceHint(src, file string, line, col1, col2 int) *CompilerEvent {
	// Create the message.
	hint := pprint.New().
		Newline().
		Yellow("  File   : " + file).Newline().
		Yellow("  Src    : " + cstr(line)).Newline().
		Yellow("  Source : " + src)

	// Create the underscore of the exact problem area.
	sx, ex := col1, col2
	if ex-sx-1 >= 0 {
		ptrLead := strings.Repeat(" ", sx-1)
		ptrDashes := strings.Repeat("-", ex-sx-1)
		var ptr string
		if ex-sx-1 == 0 {
			ptr = "           " + ptrLead + "^"
		} else {
			ptr = "           " + ptrLead + "^" + ptrDashes + "^"
		}
		hint.Newline().Yellow(ptr)
	}

	// Store the hint.
	m.hint = hint

	return m
}

func (m *CompilerEvent) Str(msg string) *CompilerEvent {
	if m.msg == nil {
		m.msg = pprint.New()
	}

	m.msg.Add(msg)
	return m
}

func (m *CompilerEvent) AddF(msg string, pairs ...string) *CompilerEvent {
	m.Str(fstr.Pairs(msg, pairs...))
	return m
}

func (m *CompilerEvent) Send() {
	// Init the output formatter.
	msg := pprint.New()

	// Prepare the level.
	lvl := m.level + " "
	switch m.level {
	case LevelInfo:
		msg.Green(lvl)
	case LevelDebug:
		msg.Gray(lvl)
	case LevelWarn:
		msg.Yellow(lvl)
	case LevelFatal, LevelError:
		msg.Red(lvl)
	default:
		msg.Green(lvl)
	}

	// Prepare the type.
	msg.Yellow("[" + m.mtype + "]: ")

	// Prepare the actual message.
	if m.msg != nil {
		msg.White(m.msg.String())
	}

	// Prepare the source hint.
	if m.hint != nil {
		msg.Newline()
		msg.Add(m.hint.String())
	}

	// Prepare the call stack.
	if m.stack != nil {
		msg.Newline()
		msg.Add(m.stack.String())
	}

	fmt.Println(msg.String())

	if m.level == LevelFatal {
		os.Exit(1)
	}
}

// Simple helper to stringify a 32-bit integer.
func cstr(i int) string {
	return strconv.Itoa(i)
}
