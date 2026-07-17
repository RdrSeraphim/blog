package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rdrseraphim/blog/internal/frontmatter"
)

func runFrontmatter(args []string) error {
	fs := flag.NewFlagSet("frontmatter", flag.ExitOnError)
	write := fs.Bool("w", false, "write the result back to the file instead of printing it to stdout")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli frontmatter [-w] <file>

Fixes up a content file's front matter to match this blog's conventions:

  - converts old TOML (+++) front matter to this blog's YAML (---) style
  - fills in a missing "lastmod" from "date"
  - fills in a missing "slug" from the file's path
  - for files under content/posts or content/pages, reorders fields to
    match the schema already used across the rest of that section

It never invents or drops a value beyond those two safe defaults, and
never touches fields it doesn't recognize - it just carries them over.
Already-clean files are left untouched (a no-op).

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

	res, err := frontmatter.Fix(path, data)
	if err != nil {
		return fmt.Errorf("%s: %w", path, err)
	}
	if !res.Changed {
		fmt.Fprintln(os.Stderr, "no front matter changes needed")
		return nil
	}

	if *write {
		if err := os.WriteFile(path, res.Body, 0o644); err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "fixed front matter in %s\n", path)
		return nil
	}

	_, err = os.Stdout.Write(res.Body)
	return err
}
