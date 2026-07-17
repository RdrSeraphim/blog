package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/rdrseraphim/blog/internal/frontmatter"
)

func runFrontmatter(args []string) error {
	fs := flag.NewFlagSet("frontmatter", flag.ExitOnError)
	write := fs.Bool("w", false, "write results back to the files instead of printing to stdout")
	yaml := fs.Bool("yaml", false, "normalize to YAML (---) instead of the canonical TOML (+++)")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli frontmatter [-w] [--yaml] [<file|dir>...]

Normalizes content front matter to this blog's conventions:

  - re-serializes to the canonical format (TOML +++ by default; --yaml
    to go the other way)
  - fills in a missing "lastmod" from "date"
  - fills in a missing "slug" from the file's path
  - for files under content/posts or content/pages, reorders fields to
    match the schema used across the rest of that section

It never invents or drops a value beyond those two safe defaults, and
never fabricates a summary. Already-canonical files are left untouched.

With no path arguments it processes every .md under ./content; a
directory is walked recursively, a file is taken as-is. Without -w the
result is printed to stdout (only useful for a single file).
`)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	target := frontmatter.FormatTOML
	if *yaml {
		target = frontmatter.FormatYAML
	}

	paths, err := collectMarkdown(fs.Args())
	if err != nil {
		return err
	}
	if len(paths) > 1 && !*write {
		return errors.New("refusing to print multiple files to stdout; pass -w to write them in place")
	}

	changed := 0
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		res, err := frontmatter.Fix(path, data, target)
		if err != nil {
			if errors.Is(err, frontmatter.ErrNoFrontMatter) {
				// Skip files without front matter (e.g. stray includes);
				// they're not ours to rewrite.
				continue
			}
			return fmt.Errorf("%s: %w", path, err)
		}
		if !res.Changed {
			continue
		}
		changed++
		if *write {
			if err := os.WriteFile(path, res.Body, 0o644); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "fixed %s\n", path)
			continue
		}
		if _, err := os.Stdout.Write(res.Body); err != nil {
			return err
		}
	}

	if *write {
		if changed == 0 {
			fmt.Fprintln(os.Stderr, "no front matter changes needed")
		} else {
			fmt.Fprintf(os.Stderr, "fixed %d file(s)\n", changed)
		}
	} else if changed == 0 {
		fmt.Fprintln(os.Stderr, "no front matter changes needed")
	}
	return nil
}
