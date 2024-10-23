package main

var cmds = make(map[string]cmd)

type cmd interface {
	ParseArgs()
	ParseFlags()
	Exec() error
}
