package argv

import lua "github.com/yuin/gopher-lua"

func Get(l *lua.LState) chan lua.LValue {
	var (
		ch = make(chan lua.LValue)
	)

	go func() {
		defer close(ch)

		for i := -l.GetTop(); i < 0; i++ {
			ch <- l.Get(i)
		}
	}()

	return ch
}
