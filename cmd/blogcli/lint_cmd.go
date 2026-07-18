package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/rdrseraphim/blog/internal/frontmatter"
)

func runLint(args []string) error {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	yaml := fs.Bool("yaml", false, "treat YAML (---) as the canonical format instead of TOML (+++)")
	warnAsError := fs.Bool("strict", false, "exit non-zero on warnings too, not just errors")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli lint [--yaml] [--strict] [<file|dir>...]

Checks content front matter against this blog's conventions without
changing anything, and reports issues grouped by file. Errors include a
wrong (non-canonical) format, unparseable front matter, a missing
required field, or an invalid date; warnings include a missing/empty
summary (which leaves the meta description uncurated), a cover without a
cover-alt, and a slug that doesn't match the file's path.

With no path arguments it checks every .md under ./content. Exits
non-zero if any errors are found (or any warnings, with --strict), so it
can gate a pre-commit hook or CI.
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

	var errorCount, warnCount int
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if _, _, _, ok := frontmatter.SplitAny(data); !ok {
			// No front matter at all - not a content file we lint.
			continue
		}
		issues := frontmatter.Lint(path, data, target)
		if len(issues) == 0 {
			continue
		}
		fmt.Printf("%s\n", path)
		for _, i := range issues {
			fmt.Printf("  %-7s %s\n", i.Severity, i.Message)
			if i.Severity == frontmatter.Error {
				errorCount++
			} else {
				warnCount++
			}
		}
	}

	fmt.Fprintf(os.Stderr, "%d error(s), %d warning(s) across %d file(s)\n", errorCount, warnCount, len(paths))
	if errorCount > 0 || (*warnAsError && warnCount > 0) {
		return errors.New("lint failed")
	}
	return nil
}
