package main

import (
	"fmt"
	"os"
	"time"

	"github.com/illbjorn/skal/internal/skal"
)

func init() {
	cmd := new(cmdExec)
	cmds["exec"] = cmd
	cmds["e"] = cmd
}

var _ cmd = &cmdExec{}

type cmdExec struct {
	input string
}

func (cmd *cmdExec) ParseArgs() {
	// Expect 2-3 args.
	switch len(os.Args) {
	case 3: // Just the input path.
		cmd.input = os.Args[2]

	case 4: // Input and output paths.
		cmd.input = os.Args[2]

	default:
		for i, arg := range os.Args {
			println(i, ":", arg)
		}
		println(len(os.Args))
		println(helpText)
		os.Exit(1)
	}
}

func (cmd *cmdExec) ParseFlags() {}

func (cmd *cmdExec) Exec() error {
	fmt.Println("Compiling and running:", cmd.input)
	start := time.Now()
	skal.CompileAndRun(cmd.input)
	dur := time.Since(start)

	println("Runtime:", dur.String())

	return nil
}
