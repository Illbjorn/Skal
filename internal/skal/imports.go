package skal

import (
	"bytes"
	"errors"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/illbjorn/fstr"
	"github.com/illbjorn/skal/internal/skal/lex/token"
	"github.com/illbjorn/skal/internal/skal/sklog"
)

func getImports(j *job) *job {
	// Identify the root directory.
	root := filepath.Dir(j.Main.Path)

	// Start the recursive import collection.
	return parseImports(root, j.Main.Content, j)
}

func parseImports(root string, fileContent string, j *job) *job {
	// Read first few lines, look for `import` keyword.
	// If we find a keyword aside from `import`, eject.
	buf := bytes.NewBufferString(fileContent)

	for {
		l, err := buf.ReadBytes('\n')
		if err != nil {
			// On an EOF, we still have data from the read. Read the last line before
			// returning.
			if errors.Is(err, io.EOF) {
				parseImportGuard(root, l, j)
				return j
			}

			sklog.CFatalF(
				"Failed to read next line of main import with error: {err}.",
				err.Error(),
			)
		}

		// Strip whitespace.
		l = bytes.TrimSpace(l)

		// Ignore comment and blank lines.
		if bytes.HasPrefix(l, []byte(token.Comment.String())) || len(l) == 0 {
			continue
		}

		// If it's not an import line, don't bother reading the rest of the File.
		if !bytes.HasPrefix(l, []byte(token.Import.String())) {
			return j
		}

		// Parse the import line.
		parseImportGuard(root, l, j)
	}
}

func parseImportGuard(root string, l []byte, j *job) *job {
	// Strip whitespace.
	l = bytes.TrimSpace(l)

	// Ignore comment and blank lines.
	if bytes.HasPrefix(l, []byte(token.Comment.String())) || len(l) == 0 {
		return j
	}

	// If it's not an import line, don't bother reading the rest of the File.
	if !bytes.HasPrefix(l, []byte(token.Import.String())) {
		return j
	}

	return parseImport(root, getImportPath(root, l), j)
}

func parseImport(root, importPath string, j *job) *job {
	// Skip duplicate imports.
	for _, i := range j.Imports {
		if i.Path == importPath {
			return j
		}
	}

	// Read the imported module.
	c, err := os.ReadFile(importPath)
	if err != nil {
		sklog.CFatalF(
			"Failed to read imported module: {path} with error: {err}.",
			"path", importPath,
			"err", err.Error(),
		)
	}

	// Recursively Evaluate for Imports.
	parseImports(root, string(c), j)

	// Append the import.
	j.Imports = append(
		j.Imports,
		&srcFile{Path: importPath, Content: string(c), Import: true})

	return j
}

var pSQImport = regexp.MustCompile("(?:'(.+?)')")
var pDQImport = regexp.MustCompile("(?:\"(.+?)\")")

func getImportPath(root string, l []byte) string {
	// Parse import module name inside quotes.
	var imp []byte
	if pSQImport.Match(l) {
		imp = pSQImport.FindSubmatch(l)[1]
	}

	if pDQImport.Match(l) {
		imp = pDQImport.FindSubmatch(l)[1]
	}

	// Split the path and rejoin (cross-OS File pathing support /\).
	paths := strings.Split(string(imp), "/")

	// Prepend the root path.
	paths = append([]string{root}, paths...)

	// Rejoin on OS sep, this will serve as the checked directory path.
	dpath := filepath.Join(paths...)

	// Append the `.sk` extension, this will serve as the checked File path.
	fpath := dpath + ".sk"

	// Identify if this path points to a directory or a File.
	// If it's a File, simply process the File.
	fstat, err := os.Stat(fpath)
	if err == nil && !fstat.IsDir() {
		return fpath
	}

	// If it's a directory, look for a `Mod.sk` File in that directory.
	dirstat, err := os.Stat(dpath)
	if err == nil && dirstat.IsDir() {
		// Look for the `Mod.sk` File.
		path := filepath.Join(dpath, "Mod.sk")
		if fstat, err = os.Stat(path); err == nil && !fstat.IsDir() {
			// The `Mod.sk` File exists.
			return path
		}
	}

	e := fstr.Pairs(
		"Failed to identify a valid import for: '{path}'.",
		"path", string(imp),
	)
	panic(e)
}
