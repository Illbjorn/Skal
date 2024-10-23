package emit

import (
	"strings"
)

// TODO: Eliminate package level state!
var stack = Stack{i: 1, s: make([][]string, 10)}

type Stack struct {
	s             [][]string
	i             int
	defersTotal   int
	defersInScope int
}

func (s *Stack) Indent() string {
	return strings.Repeat(" ", s.i*2)
}

func (s *Stack) Add(deferral string) {
	s.defersInScope++
	s.defersTotal++
	s.s[s.i] = append([]string{deferral}, s.s[s.i]...)
}

func (s *Stack) Push() {
	s.i++
	s.defersInScope = 0
}

func (s *Stack) Pop() {
	s.s[s.i] = make([]string, 0)
	if s.i > 0 {
		s.i--
	}

	s.defersTotal -= s.defersInScope
	s.defersInScope = len(s.s[s.i])
}

type deferral struct {
	stmt    string
	current int
	total   int
}

func (s *Stack) Unwind() chan deferral {
	ch := make(chan deferral)

	go func() {
		defer close(ch)

		// Ship deferrals.
		current := 0
		total := s.defersTotal
		for i := len(s.s) - 1; i > 0; i-- {
			if len(s.s[i]) == 0 {
				continue
			}

			for _, v := range s.s[i] {
				ch <- deferral{v, current, total}
				current++
			}
		}
	}()

	return ch
}
