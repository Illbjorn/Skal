package skal

import (
	"os"

	"github.com/illbjorn/skal/internal/skal/emit"
	"github.com/illbjorn/skal/internal/skal/lex"
	"github.com/illbjorn/skal/internal/skal/parse"
	"github.com/illbjorn/skal/internal/skal/sklog"
	"github.com/illbjorn/skal/internal/skal/typeset"
	"github.com/illbjorn/skal/pkg/formatter"
)

func Compile(inputPath, outputPath string) {
	// Entrypoint I/O
	// Read the 'main' File.
	b, err := os.ReadFile(inputPath)
	if err != nil {
		sklog.CFatalF(
			"Failed to read input File: '{path}', with error: {err}.",
			"path", inputPath,
			"err", err.Error(),
		)
	}

	// Open a writable stream to the output File.
	outFile, err := os.OpenFile(
		outputPath,
		os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
		0600)
	if err != nil {
		sklog.CFatalF(
			"Failed to open writable stream to output File with error: {err}.",
			"err", err.Error(),
		)
	}
	defer func(outFile *os.File) {
		if err := outFile.Close(); err != nil {
			println("ERROR: Failed to close output file", outFile.Name()+".")
			println("Inner error:", err.Error()+".")
		}
	}(outFile)

	// Init the job.
	j := &job{
		Main: &srcFile{
			Path:    inputPath,
			Content: string(b),
		},
		OutputPath: outputPath,
	}

	// Assemble all imported modules.
	j = getImports(j)

	// Write the basic env header.
	_, _ = outFile.Write(tmplHeader)

	// Init the formatter we'll use to produce the formatted output.
	// We'll reset this underlying buffer between files to avoid some unnecessary
	// reallocation.
	emt := formatter.NewFormatter()

	// Process source files.
	for inFile := range j.Gen() {
		compileFile(inFile, outFile, emt)
	}

	// Write the basic env footer.
	_, _ = outFile.Write(tmplFooter)
}

// Compile a single provided File.
// Lex -> Parse -> Typeset -> Emit
func compileFile(inFile *srcFile, outFile *os.File, emt *formatter.Formatter) {
	// Ignore empty files.
	if len(inFile.Content) == 0 {
		return
	}

	// Reset buffers.
	defer emt.Reset()

	// Lex
	tokenCollection := lex.Lex(inFile.Path, string(inFile.Content))

	// Parse
	tree := parse.Parse(tokenCollection, nil, nil)

	// Typeset
	set := typeset.Typeset(tree)

	// Emit
	compiled := emit.Emit(set, inFile.Path, inFile.Import, emt)

	// Write the compiled code.
	_, _ = outFile.Write(compiled)
}

var (
	// The boilerplate File header which sets up the script's executing environment.
	tmplHeader = []byte(`-- Create the app environment.
local __ENV__ = { __index = _G }
-- Set the metatable.
setmetatable(__ENV__, __ENV__)
-- DEBUG: Export the app environment table.
_G._ENV_ = __ENV__
-- Open the app function.
local function __LOAD__()`)

	// The boilerplate File header which closes the script executing environment
	// entrypoint and calls it.
	tmplFooter = []byte(`
end
-- Set the app function's environment to the app environment table.
setfenv(__LOAD__, __ENV__)
-- Launch the application.
__LOAD__()`)
)
