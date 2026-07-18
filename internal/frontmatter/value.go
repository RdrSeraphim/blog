package frontmatter

import (
	"regexp"
	"strings"
	"time"
)

// Format identifies a front matter serialization.
type Format int

const (
	// FormatTOML is Hugo's +++ TOML front matter, this blog's canonical style.
	FormatTOML Format = iota
	// FormatYAML is --- YAML front matter (what the Ghost-migrated content
	// originally arrived as).
	FormatYAML
)

// kind is the logical type of a front matter value, tracked so a value can
// be re-serialized correctly into either format regardless of which one it
// was read from.
type kind int

const (
	kindString kind = iota
	kindBool
	kindDateTime
	kindArray
)

// Value is a decoded front matter value plus enough type information to
// round-trip it into TOML or YAML. Datetimes are kept as their original
// opaque token (never parsed and reformatted) so exact precision/offset is
// preserved.
type Value struct {
	kind  kind
	str   string   // string text (decoded), datetime token, or bool literal
	items []string // element strings for kindArray
}

var rfc3339ish = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}[Tt ]\d{2}:\d{2}:\d{2}(\.\d+)?([Zz]|[+-]\d{2}:\d{2})?$`)

func looksLikeDateTime(raw string) bool {
	if !rfc3339ish.MatchString(raw) {
		return false
	}
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05"} {
		if _, err := time.Parse(layout, raw); err == nil {
			return true
		}
	}
	return false
}

// DecodeValue interprets the raw text after a "key:" / "key =" separator,
// in the given source format, into a typed Value.
func DecodeValue(raw string, src Format) Value {
	raw = strings.TrimSpace(raw)

	if raw == "true" || raw == "false" {
		return Value{kind: kindBool, str: raw}
	}
	if strings.HasPrefix(raw, "[") && strings.HasSuffix(raw, "]") {
		return Value{kind: kindArray, items: parseArray(raw, src)}
	}
	if looksLikeDateTime(raw) {
		return Value{kind: kindDateTime, str: raw}
	}
	return Value{kind: kindString, str: decodeScalarString(raw, src)}
}

// decodeScalarString unwraps a single scalar string value written in the
// source format, yielding its literal text.
func decodeScalarString(raw string, src Format) string {
	if len(raw) >= 2 {
		q := raw[0]
		if q == raw[len(raw)-1] {
			inner := raw[1 : len(raw)-1]
			switch {
			case q == '\'' && src == FormatYAML:
				return strings.ReplaceAll(inner, "''", "'")
			case q == '\'' && src == FormatTOML:
				return inner // TOML literal string: no escaping
			case q == '"':
				return unescapeBackslash(inner)
			}
		}
	}
	return raw
}

// unescapeBackslash processes the backslash escapes common to YAML
// double-quoted and TOML basic strings.
func unescapeBackslash(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != '\\' || i+1 >= len(s) {
			b.WriteByte(s[i])
			continue
		}
		i++
		switch s[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		case 'r':
			b.WriteByte('\r')
		case '"':
			b.WriteByte('"')
		case '\\':
			b.WriteByte('\\')
		default:
			b.WriteByte('\\')
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// parseArray splits a bracketed list into its decoded element strings,
// respecting quotes so commas inside an element don't split it.
func parseArray(raw string, src Format) []string {
	inner := strings.TrimSpace(raw[1 : len(raw)-1])
	if inner == "" {
		return nil
	}
	var items []string
	var cur strings.Builder
	var quote byte
	for i := 0; i < len(inner); i++ {
		c := inner[i]
		if quote != 0 {
			cur.WriteByte(c)
			if c == quote {
				quote = 0
			}
			continue
		}
		switch c {
		case '\'', '"':
			quote = c
			cur.WriteByte(c)
		case ',':
			items = append(items, decodeScalarString(strings.TrimSpace(cur.String()), src))
			cur.Reset()
		default:
			cur.WriteByte(c)
		}
	}
	if strings.TrimSpace(cur.String()) != "" {
		items = append(items, decodeScalarString(strings.TrimSpace(cur.String()), src))
	}
	return items
}

// Encode serializes v into the given format's value syntax.
func (v Value) Encode(f Format) string {
	switch v.kind {
	case kindBool, kindDateTime:
		return v.str
	case kindArray:
		return encodeArray(v.items, f)
	default:
		return encodeString(v.str, f)
	}
}

func encodeArray(items []string, f Format) string {
	if len(items) == 0 {
		return "[]"
	}
	parts := make([]string, len(items))
	for i, it := range items {
		parts[i] = encodeString(it, f)
	}
	sep := ", "
	if f == FormatYAML {
		sep = ","
	}
	return "[" + strings.Join(parts, sep) + "]"
}

func hasControl(s string) bool {
	for _, r := range s {
		if r < 0x20 {
			return true
		}
	}
	return false
}

var safeYAMLPlain = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9 ]*$`)

func encodeString(s string, f Format) string {
	if f == FormatTOML {
		// Prefer a literal string (no escaping); fall back to a basic
		// string only when the text contains a single quote or a control
		// character a literal string can't hold.
		if !strings.ContainsRune(s, '\'') && !hasControl(s) {
			return "'" + s + "'"
		}
		return `"` + escapeBasic(s) + `"`
	}

	// YAML: a simple word/phrase can go unquoted; anything with quotes or
	// punctuation that YAML would misread gets a single-quoted (or, if it
	// contains a single quote, double-quoted) string.
	if s != "" && safeYAMLPlain.MatchString(s) && !looksLikeDateTime(s) && s != "true" && s != "false" {
		return s
	}
	if !strings.ContainsRune(s, '\'') && !hasControl(s) {
		return "'" + s + "'"
	}
	return `"` + escapeBasic(s) + `"`
}

func escapeBasic(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '"':
			b.WriteString(`\"`)
		case '\\':
			b.WriteString(`\\`)
		case '\n':
			b.WriteString(`\n`)
		case '\t':
			b.WriteString(`\t`)
		case '\r':
			b.WriteString(`\r`)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}
