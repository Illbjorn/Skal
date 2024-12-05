package main

import (
	"os"
	"strings"

	"github.com/illbjorn/fstr"
)

var helpText = fstr.Pairs(`
{magenta}The Skal compiler and runtime.{reset}

{cyan}Docs{reset} : https://skal.dev/docs
{cyan}Bugs{reset} : https://github.com/Illbjorn/Skal/issues

To compile a script:

  {green}skal c ./main.sk{reset}

Commands:
  compile, c       Compile a Skal script.

Options:
  --with-perf, -p  Produce performance measurement output after a compile.
	--dump-ast,  -d  Serialize the built AST to JSON and write to file.
	--watch,     -w  Watch the targeted source file and recompile on change.
	--help,      -h  How you got here!
`,
	"cyan", "\033[0m",
	"yellow", "\033[33m",
	"magenta", "\033[35m",
	"reset", "\033[0m",
	"green", "\033[32m",
)

func main() {
	if len(os.Args) < 2 {
		printHelpTextAndExit()
	}

	// Look for the called command.
	called := strings.ToLower(os.Args[1])
	if called, ok := cmds[called]; !ok {
		printHelpTextAndExit()
	} else {
		called.ParseArgs()
		called.ParseFlags()
		if err := called.Exec(); err != nil {
			printHelpTextAndExit()
		}
	}
}
