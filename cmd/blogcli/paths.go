package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// collectMarkdown expands the given path arguments into a list of markdown
// files: a directory is walked recursively for *.md, a file is taken
// as-is. With no args, it walks ./content. This is what lets the frontmatter
// and lint commands operate on a single file or the whole tree.
func collectMarkdown(args []string) ([]string, error) {
	if len(args) == 0 {
		args = []string{"content"}
	}
	var out []string
	seen := map[string]bool{}
	add := func(p string) {
		if !seen[p] {
			seen[p] = true
			out = append(out, p)
		}
	}

	for _, arg := range args {
		info, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			add(filepath.Clean(arg))
			continue
		}
		err = filepath.WalkDir(arg, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() && strings.HasSuffix(p, ".md") {
				add(filepath.Clean(p))
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
