package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rdrseraphim/blog/internal/footnotes"
	"github.com/rdrseraphim/blog/internal/frontmatter"
)

func runFootnotes(args []string) error {
	fs := flag.NewFlagSet("footnotes", flag.ExitOnError)
	write := fs.Bool("w", false, "write the result back to the file instead of printing it to stdout")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli footnotes [-w] <file>

Renumbers mnemonic footnotes, e.g.:

    Point one.[^meow: A cat sound.] Point two, mentioned again[^meow].

into the sequential style used across the blog:

    Point one.[^1] Point two, mentioned again[^1].

    [^1]: A cat sound.

Bare references paired with a "[^label]: text" line elsewhere in the
document work too, so re-running this on an already-numbered post is a
no-op.
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
		return fmt.Errorf("no definition found for: [^%s] (add a \"[^%s]: ...\" line, or use the [^%s: text] inline shorthand)",
			strings.Join(res.Missing, "], [^"), res.Missing[0], res.Missing[0])
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

	if *write {
		if err := os.WriteFile(path, out, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "renumbered footnotes in %s\n", path)
		return nil
	}

	_, err = os.Stdout.Write(out)
	return err
}
