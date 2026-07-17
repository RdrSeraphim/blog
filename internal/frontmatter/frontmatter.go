// Package frontmatter splits Hugo content files into their front matter
// block and Markdown body, and formats posts back into that shape.
package frontmatter

import (
	"fmt"
	"regexp"
	"strings"
)

const delim = "---"

// Split separates data into its front matter (without the delimiters) and
// body. ok is false if the file has no front matter block, in which case
// body is the whole input.
func Split(data []byte) (fm string, body string, ok bool) {
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

// Join reassembles a front matter block and body into a full file.
func Join(fm string, body string) []byte {
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

const tomlDelim = "+++"

// SplitAny is like Split, but also recognizes TOML (+++ fenced) front
// matter, which Hugo accepts equally but this blog otherwise never uses.
func SplitAny(data []byte) (fm string, body string, isTOML bool, ok bool) {
	if fm, body, ok := Split(data); ok {
		return fm, body, false, true
	}
	s := string(data)
	if !strings.HasPrefix(s, tomlDelim+"\n") {
		return "", s, false, false
	}
	rest := s[len(tomlDelim)+1:]
	idx := strings.Index(rest, "\n"+tomlDelim)
	if idx == -1 {
		return "", s, false, false
	}
	fm = rest[:idx]
	after := rest[idx+len(tomlDelim)+1:]
	body = strings.TrimPrefix(after, "\n")
	return fm, body, true, true
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

var safePlainScalar = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 _-]*$`)
var looksLikeDate = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[Tt]\d{2}:\d{2}:\d{2}`)

// Unquote strips a value's surrounding quotes if doing so still leaves a
// safe, unambiguous plain YAML scalar (a simple word/sentence-ish string,
// or an RFC 3339-ish date) - used to clean up values that only needed
// quoting because of TOML's stricter string syntax. Anything else (values
// with colons, leading punctuation, etc., where quoting is load-bearing)
// is returned unchanged.
func Unquote(value string) string {
	if len(value) < 2 {
		return value
	}
	quote := value[0]
	if (quote != '\'' && quote != '"') || value[len(value)-1] != quote {
		return value
	}
	inner := value[1 : len(value)-1]
	if safePlainScalar.MatchString(inner) || looksLikeDate.MatchString(inner) {
		return inner
	}
	return value
}

// RenderFields writes fields back out as a YAML front matter block.
func RenderFields(fields []Field) string {
	var b strings.Builder
	for _, f := range fields {
		fmt.Fprintf(&b, "%s: %s\n", f.Key, f.Value)
	}
	return b.String()
}
