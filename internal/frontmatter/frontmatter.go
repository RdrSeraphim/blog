// Package frontmatter splits Hugo content files into their front matter
// block and Markdown body, parses the flat key/value front matter this
// blog uses, and re-serializes it into either TOML (+++) or YAML (---).
package frontmatter

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	yamlDelim = "---"
	tomlDelim = "+++"
)

// Split separates data into its YAML front matter (without the delimiters)
// and body. ok is false if the file has no YAML front matter block, in
// which case body is the whole input.
func Split(data []byte) (fm string, body string, ok bool) {
	return splitDelim(data, yamlDelim)
}

// SplitAny is like Split, but recognizes both YAML (---) and TOML (+++)
// front matter. isTOML reports which was found.
func SplitAny(data []byte) (fm string, body string, isTOML bool, ok bool) {
	if fm, body, ok := splitDelim(data, yamlDelim); ok {
		return fm, body, false, true
	}
	if fm, body, ok := splitDelim(data, tomlDelim); ok {
		return fm, body, true, true
	}
	return "", string(data), false, false
}

func splitDelim(data []byte, delim string) (fm string, body string, ok bool) {
	s := string(data)
	if !strings.HasPrefix(s, delim+"\n") {
		return "", s, false
	}
	rest := s[len(delim)+1:]
	idx := strings.Index(rest, "\n"+delim)
	if idx == -1 {
		return "", s, false
	}
	fm = rest[:idx]
	after := rest[idx+len(delim)+1:]
	body = strings.TrimPrefix(after, "\n")
	return fm, body, true
}

// Join reassembles a YAML front matter block and body into a full file.
func Join(fm string, body string) []byte {
	return JoinFormat(fm, body, FormatYAML)
}

// JoinFormat reassembles a front matter block and body into a full file,
// fencing it with the delimiter for the given format.
func JoinFormat(fm string, body string, f Format) []byte {
	delim := yamlDelim
	if f == FormatTOML {
		delim = tomlDelim
	}
	var b strings.Builder
	b.WriteString(delim)
	b.WriteString("\n")
	b.WriteString(strings.TrimRight(fm, "\n"))
	b.WriteString("\n")
	b.WriteString(delim)
	b.WriteString("\n\n")
	b.WriteString(strings.TrimLeft(body, "\n"))
	return []byte(b.String())
}

// Field is a single front matter key and the raw text of its value,
// exactly as written (still quoted, bracketed, etc. if it was).
type Field struct {
	Key   string
	Value string
}

var fieldLineRe = regexp.MustCompile(`^([A-Za-z_][A-Za-z0-9_-]*)\s*[:=]\s?(.*)$`)

// ParseFields reads a flat front matter block (TOML or YAML, no nested
// tables/mappings - which is all this blog's content ever uses) into an
// ordered list of key/value pairs. It errors on anything it can't
// confidently read as a single "key: value" or "key = value" line, rather
// than silently dropping content.
func ParseFields(raw string) ([]Field, error) {
	var fields []Field
	for _, line := range strings.Split(raw, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		m := fieldLineRe.FindStringSubmatch(line)
		if m == nil {
			return nil, fmt.Errorf("can't parse front matter line: %q", line)
		}
		fields = append(fields, Field{Key: m[1], Value: m[2]})
	}
	return fields, nil
}
