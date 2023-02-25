package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/viegasfh/go-vimlc"
	"github.com/viegasfh/go-vimlc/compiler/sexpression"
)

var sexpr = flag.Bool("sexpr", false, "output s-expression")
var usejson = flag.Bool("json", false, "output json")

func main() {
	flag.Parse()

	opt := &vimlparser.ParseOption{Neovim: *sexpr}

	if len(flag.Args()) == 0 {
		if err := parseFile("", os.Stdin, os.Stdout, opt, *usejson); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	exitCode := 0

	for i, file := range flag.Args() {
		f, err := os.Open(file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if i == 0 {
				flag.Usage()
			}
			exitCode = 1
			continue
		}
		if err := parseFile(f.Name(), f, os.Stdout, opt, *usejson); err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}

// filename is empty string if r is os.Stdin
func parseFile(filename string, r io.ReadCloser, w io.Writer, opt *vimlparser.ParseOption, usejson bool) error {
	defer r.Close()
	node, err := vimlparser.ParseFile(r, filename, opt)
	if err != nil {
		return err
	}
	if usejson {
		e := json.NewEncoder(w)
		e.SetIndent("", "    ")
		return e.Encode(node)
	} else {
		c := &sexpression.Compiler{Config: sexpression.Config{Indent: "  "}}
		if err := c.Compile(w, node); err != nil {
			return err
		}
	}

	return nil
}
