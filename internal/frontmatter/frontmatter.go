// Package frontmatter splits Hugo content files into their YAML front
// matter block and Markdown body, and formats posts back into that shape.
package frontmatter

import "strings"

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
