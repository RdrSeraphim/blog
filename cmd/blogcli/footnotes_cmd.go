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

	fm, body, isTOML, hasFM := frontmatter.SplitAny(data)

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
		format := frontmatter.FormatYAML
		if isTOML {
			format = frontmatter.FormatTOML
		}
		out = frontmatter.JoinFormat(fm, res.Body, format)
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
