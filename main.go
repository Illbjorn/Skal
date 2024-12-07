package main

import (
	"fmt"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	lua "github.com/yuin/gopher-lua"
)

type x struct {
	A string
	B int
	C bool
	D another
}

type another struct {
	E float64
	F uint8
}

func main() {
	l := lua.NewState()

	y := x{
		A: "Hello",
		B: 9,
		C: true,
		D: another{
			E: 1.223,
			F: 122,
		},
	}

	res := conv.StructToLTable(y, l)

	fmt.Printf("%#v\n", res)
}
