// Command blogcli handles two bits of post-authoring this blog's tooling
// otherwise doesn't cover: renumbering hand-authored footnotes, and fixing
// up front matter on content that didn't get created from a proper
// archetype (wrong format, missing fields).
//
// Metadata scaffolding for new content is handled by `hugo new` and the
// archetypes in archetypes/, and previewing by `hugo server` (or an
// in-editor renderer like render-markdown.nvim).
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
	case "footnotes":
		err = runFootnotes(os.Args[2:])
	case "frontmatter":
		err = runFrontmatter(os.Args[2:])
	case "lint":
		err = runLint(os.Args[2:])
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
  blogcli footnotes [-w] <file>            renumber footnotes to [^1], [^2]...
  blogcli frontmatter [-w] [<file|dir>...] normalize front matter (TOML)
  blogcli lint [<file|dir>...]             report front matter problems

Run 'blogcli <command> -h' for command-specific details.
`)
}
