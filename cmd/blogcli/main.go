// Command blogcli is a small CLI/TUI companion for authoring posts in
// this Hugo blog: it scaffolds new post metadata, renumbers mnemonic
// footnotes into the blog's [^1]/[^2] endnote style, and previews a post
// in the terminal or in a browser.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	var err error
	switch os.Args[1] {
	case "new":
		err = runNew(os.Args[2:])
	case "footnotes":
		err = runFootnotes(os.Args[2:])
	case "preview":
		err = runPreview(os.Args[2:])
	case "-h", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "blogcli: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "blogcli: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `blogcli - authoring companion for this Hugo blog

Usage:
  blogcli new                       interactively scaffold a new post
  blogcli footnotes [-w] <file>     renumber mnemonic footnotes to [^1], [^2]...
  blogcli preview [--html] <file>   preview a post in the terminal, or in a browser

Run 'blogcli <command> -h' for command-specific flags.
`)
}
