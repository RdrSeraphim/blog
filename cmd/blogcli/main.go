// Command blogcli renumbers hand-authored footnotes in this Hugo blog's
// posts into the sequential [^1]/[^2] endnote style used throughout the
// site.
//
// Metadata scaffolding is handled by `hugo new` and the archetypes in
// archetypes/, and previewing by `hugo server` (or an in-editor renderer
// like render-markdown.nvim) — this tool only does the one thing neither
// of those cover.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rdrseraphim/blog/internal/footnotes"
	"github.com/rdrseraphim/blog/internal/frontmatter"
)

func main() {
	flag.Usage = usage
	write := flag.Bool("w", false, "write the result back to the file instead of printing it to stdout")
	flag.Parse()

	if flag.NArg() != 1 {
		usage()
		os.Exit(2)
	}

	if err := run(flag.Arg(0), *write); err != nil {
		fmt.Fprintf(os.Stderr, "blogcli: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Usage: blogcli [-w] <file>

Renumbers footnotes written with their text right in the marker, e.g.:

    Point one.[^A cat sound.] Point two.[^A dog sound.]

into the sequential style used across the blog:

    Point one.[^1] Point two.[^2]

    [^1]: A cat sound.
    [^2]: A dog sound.

A purely numeric marker like [^1] paired with a "[^1]: text" line
elsewhere in the document is read as an already-numbered reference
rather than literal text, so re-running this on an already-numbered
post is a no-op.

Without -w, the result is printed to stdout and the file is left
untouched.
`)
}

func run(path string, write bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fm, body, hasFM := frontmatter.Split(data)

	res, err := footnotes.Renumber(body)
	if err != nil {
		return err
	}
	if len(res.Missing) > 0 {
		return fmt.Errorf("no definition found for numeric reference(s) [^%s] (add a \"[^%s]: ...\" line, or replace it with the note text directly, e.g. [^the note text])",
			strings.Join(res.Missing, "], [^"), res.Missing[0])
	}
	if !res.Changed {
		fmt.Fprintln(os.Stderr, "no footnote changes needed")
		return nil
	}

	var out []byte
	if hasFM {
		out = frontmatter.Join(fm, res.Body)
	} else {
		out = []byte(res.Body)
	}

	if write {
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "renumbered footnotes in %s\n", path)
		return nil
	}

	_, err = os.Stdout.Write(out)
	return err
}
