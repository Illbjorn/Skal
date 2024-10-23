package pprint

import (
	"bytes"
	"fmt"
)

func New() *Msg {
	return &Msg{
		buffer: bytes.NewBuffer(nil),
	}
}

type Msg struct {
	buffer *bytes.Buffer
}

func (c Msg) String() string {
	return c.buffer.String()
}

func (c *Msg) Magenta(msg string) *Msg {
	c.buffer.WriteString(string(colorMagenta) + msg + string(colorReset))
	return c
}

func (c *Msg) Cyan(msg string) *Msg {
	c.buffer.WriteString(string(colorCyan) + msg + string(colorReset))
	return c
}

func (c *Msg) Gray(msg string) *Msg {
	c.buffer.WriteString(string(colorGray) + msg + string(colorReset))
	return c
}

func (c *Msg) Red(msg string) *Msg {
	c.buffer.WriteString(string(colorRed) + msg + string(colorReset))
	return c
}

func (c *Msg) Yellow(msg string) *Msg {
	c.buffer.WriteString(string(colorYellow) + msg + string(colorReset))
	return c
}

func (c *Msg) White(msg string) *Msg {
	c.buffer.WriteString(string(colorWhite) + msg + string(colorReset))
	return c
}

func (c *Msg) Green(msg string) *Msg {
	c.buffer.WriteString(string(colorGreen) + msg + string(colorReset))
	return c
}

func (c *Msg) Blue(msg string) *Msg {
	c.buffer.WriteString(string(colorBlue) + msg + string(colorReset))
	return c
}

func (c *Msg) Add(msg string) *Msg {
	c.buffer.WriteString(msg)
	return c
}

func (c *Msg) Newline() *Msg {
	return c.Add("\n")
}

func (c *Msg) Println() {
	fmt.Println(c.String())
}
