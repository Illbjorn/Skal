package main

import (
	"bytes"
	"crypto"
	"encoding/binary"
	"flag"
	"io/fs"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"time"

	"github.com/illbjorn/skal/internal/skal"
)

func init() {
	cmd := &cmdCompile{}
	cmds["compile"] = cmd
	cmds["c"] = cmd
}

var _ cmd = &cmdCompile{}

type cmdCompile struct {
	input  string
	output string
	watch  bool
}

func (cmd *cmdCompile) ParseArgs() {
	// Expect 2-3 args.
	switch len(os.Args) {
	case 3: // Just the input path.
		cmd.input = os.Args[2]
		cmd.output = strings.Replace(cmd.input, filepath.Ext(cmd.input), ".lua", 1)

	case 4: // Input and output paths.
		cmd.input = os.Args[2]
		cmd.output = os.Args[3]

	default:
		for i, arg := range os.Args {
			println(i, ":", arg)
		}
		println(len(os.Args))
		printHelpTextAndExit()
	}
}

func (cmd *cmdCompile) ParseFlags() {
	//--
	// Define flags
	watch := flag.Bool("watch", false, "")

	// Parse
	flag.Parse()

	//--
	// Assign flags
	cmd.watch = *watch
}

func (cmd *cmdCompile) Exec() error {
	// Compile!
	if cmd.watch {
		watchCompile(cmd.input, cmd.output)
	} else {
		compile(cmd.input, cmd.output)
	}

	return nil
}

func watchCompile(input, output string) {
	done := make(chan os.Signal, 2)
	signal.Notify(done, os.Interrupt)

	var h []byte
	for {
		select {
		case <-done:
			return
		case <-time.After(200 * time.Millisecond):
			nh := hash(input)
			if !bytes.Equal(nh, h) {
				compile(input, output)
				h = nh
			}
		}
	}
}

func compile(input, output string) {
	start := time.Now()
	skal.Compile(input, output)
	dur := time.Since(start)
	println("Compile Time:", dur.String())
}

// TODO: Currently this just walks the source directory and hashes all .sk files
// We should eventually instead actually parse the provided source file,
// recursively enumerate all imports and watch those.
func hash(input string) []byte {
	d := crypto.MD5.New()

	// Get the directory of the input file.
	dir := filepath.Dir(input)

	err := filepath.WalkDir(
		dir, func(path string, item fs.DirEntry, _ error) error {
			// Rule out non-Skal files.
			if item.IsDir() || filepath.Ext(item.Name()) != ".sk" {
				return nil
			}

			i, err := item.Info()
			if err != nil {
				return err
			}

			// File Name
			fn := []byte(i.Name())
			d.Write(fn)

			// File Mod Time
			mt := []byte(i.ModTime().String())
			d.Write(mt)

			// File Size
			fsz := i.Size()
			b := make([]byte, 8)
			binary.LittleEndian.PutUint64(b, uint64(fsz))
			d.Write(b)

			return nil
		})
	if err != nil {
		println("ERROR: Failed to calculate source tree hash. Inner error:", err.Error()+".")
	}

	return d.Sum(nil)
}
