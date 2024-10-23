package formatter

import (
	"bytes"
)

func NewFormatter() *Formatter {
	return &Formatter{buf: bytes.NewBuffer(nil)}
}

type Formatter struct {
	buf   *bytes.Buffer
	hooks []FormatterHook
}

// Describes a function signature for use with Formatter.Hook().
type FormatterHook func(string) string

// Set a hook to mutate future input values until the returned `unset` callback
// is called.
func (f *Formatter) Hook(fn FormatterHook) func() {
	// Add the callback.
	f.hooks = append(f.hooks, fn)

	index := len(f.hooks) - 1

	// Return a callback to unset the hook.
	return func() {
		f.hooks[index] = f.hooks[len(f.hooks)-1]
		f.hooks = f.hooks[:len(f.hooks)-1]
	}
}

// Reset the underlying buffer.
func (f *Formatter) Reset() *Formatter {
	f.buf.Reset()
	return f
}

// Write a string value to the buffer.
func (f *Formatter) Str(v string) *Formatter {
	for _, h := range f.hooks {
		v = h(v)
	}

	if len(v) == 0 {
		return f
	}

	f.buf.WriteString(v)
	return f
}

// Write a white space to the underlying buffer.
func (f *Formatter) Space() *Formatter {
	f.buf.WriteString(" ")
	return f
}

// Write a newline character to the buffer.
func (f *Formatter) Newline() *Formatter {
	f.buf.WriteString("\n")
	return f
}

// Stringify the underlying buffer.
func (f *Formatter) String() string {
	return f.buf.String()
}

// Return the underlying buffer's byte slice.
func (f *Formatter) Bytes() []byte {
	return f.buf.Bytes()
}
