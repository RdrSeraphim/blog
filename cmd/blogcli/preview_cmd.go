package main

import (
	"flag"
	"fmt"
	"os"
)

func runPreview(args []string) error {
	fs := flag.NewFlagSet("preview", flag.ExitOnError)
	htmlMode := fs.Bool("html", false, "serve a live-reloading HTML preview in a browser instead of the terminal TUI")
	noOpen := fs.Bool("no-open", false, "with --html, print the URL instead of trying to open a browser")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli preview [--html] [--no-open] <file>

Without flags, opens a terminal viewer (glamour-rendered, scrollable,
auto-reloads when the file changes on disk).

With --html, starts a local HTTP server rendering the post to styled,
theme-aware HTML and opens it in a browser; it also auto-reloads on save.
`)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		os.Exit(2)
	}
	path := fs.Arg(0)

	if _, err := os.Stat(path); err != nil {
		return err
	}

	if *htmlMode {
		return previewHTML(path, *noOpen)
	}
	return previewTUI(path)
}
