// Package footnotes renumbers hand-authored, mnemonic footnotes (e.g.
// [^meow] or the inline shorthand [^meow: text]) into the sequential
// [^1], [^2], ... style used throughout the blog, collecting their
// definitions into an endnote block at the bottom of the document.
package footnotes

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// Result is the outcome of a Renumber call.
type Result struct {
	Body    string   // the rewritten body (unchanged if nothing to do)
	Changed bool     // whether Body differs from the input
	Missing []string // labels referenced but never given definition text
}

var trailingDefRe = regexp.MustCompile(`^\[\^([^\]\s]+)\]:[ \t]?(.*)$`)

// Renumber rewrites all footnote references in body to sequential numeric
// identifiers, in order of first appearance, and appends a clean endnote
// block containing their definitions. It supports three ways of supplying
// definition text, and they may be mixed freely within one document:
//
//   - inline shorthand at the point of use: [^label: definition text]
//   - a bare reference [^label] paired with a line of the form
//     "[^label]: definition text" anywhere else in the document (this is
//     also how already-numbered posts like [^1] ... [^1]: ... are read
//     back in, so re-running Renumber on an already-clean post is a no-op)
//
// If any referenced label has no definition anywhere in the document, it
// is left untouched in the body and reported in Missing; callers should
// treat a non-empty Missing as reason not to write the result out.
func Renumber(body string) (Result, error) {
	content, trailingDefs := extractTrailingDefs(body)

	defText := map[string]string{}
	for label, text := range trailingDefs {
		defText[label] = text
	}

	var out strings.Builder
	order := []string{}
	numOf := map[string]int{}
	missingSet := map[string]bool{}

	i := 0
	for i < len(content) {
		start := strings.Index(content[i:], "[^")
		if start == -1 {
			out.WriteString(content[i:])
			break
		}
		start += i
		out.WriteString(content[i:start])

		label, defTextInline, hasInlineDef, tokenEnd, isToken := scanToken(content, start)
		if !isToken {
			out.WriteString("[^")
			i = start + 2
			continue
		}

		if hasInlineDef {
			defText[label] = defTextInline
		}

		if _, seen := numOf[label]; !seen {
			numOf[label] = len(order) + 1
			order = append(order, label)
		}

		if _, ok := defText[label]; !ok {
			missingSet[label] = true
			// Leave the original token untouched so the document isn't
			// silently corrupted; caller decides whether to abort.
			out.WriteString(content[start:tokenEnd])
		} else {
			fmt.Fprintf(&out, "[^%d]", numOf[label])
		}

		i = tokenEnd
	}

	if len(order) == 0 && len(trailingDefs) == 0 {
		return Result{Body: body, Changed: false}, nil
	}

	missing := make([]string, 0, len(missingSet))
	for l := range missingSet {
		missing = append(missing, l)
	}
	sort.Strings(missing)
	if len(missing) > 0 {
		return Result{Body: body, Changed: false, Missing: missing}, nil
	}

	newBody := strings.TrimRight(out.String(), "\n") + "\n"
	if len(order) > 0 {
		var defs strings.Builder
		for _, label := range order {
			fmt.Fprintf(&defs, "[^%d]: %s\n", numOf[label], defText[label])
		}
		newBody += "\n" + strings.TrimRight(defs.String(), "\n") + "\n"
	}

	changed := newBody != body
	return Result{Body: newBody, Changed: changed, Missing: nil}, nil
}

// extractTrailingDefs strips a contiguous block of "[^label]: text" lines
// from the end of body (skipping a single blank-line separator, matching
// how existing posts already look) and returns the remaining content plus
// the stripped definitions.
func extractTrailingDefs(body string) (content string, defs map[string]string) {
	lines := strings.Split(body, "\n")
	defs = map[string]string{}

	end := len(lines)
	for end > 0 && strings.TrimSpace(lines[end-1]) == "" {
		end--
	}

	start := end
	for start > 0 {
		m := trailingDefRe.FindStringSubmatch(lines[start-1])
		if m == nil {
			break
		}
		defs[m[1]] = m[2]
		start--
	}

	if start == end {
		// no trailing def block found
		return body, defs
	}

	remaining := lines[:start]
	for len(remaining) > 0 && strings.TrimSpace(remaining[len(remaining)-1]) == "" {
		remaining = remaining[:len(remaining)-1]
	}
	return strings.Join(remaining, "\n") + "\n", defs
}

// scanToken attempts to parse a footnote token starting at content[start]
// (which must be "[^"). It returns the label, an inline definition if the
// token used the "[^label: text]" shorthand, and the index just past the
// token's closing bracket. isToken is false if content[start:] isn't a
// well-formed footnote token, in which case only the leading "[^" should
// be consumed by the caller.
func scanToken(content string, start int) (label string, defText string, hasInlineDef bool, end int, isToken bool) {
	i := start + 2
	labelStart := i
	for i < len(content) && isLabelChar(content[i]) {
		i++
	}
	if i == labelStart {
		return "", "", false, 0, false
	}
	label = content[labelStart:i]

	if i < len(content) && content[i] == ']' {
		return label, "", false, i + 1, true
	}

	if i < len(content) && content[i] == ':' {
		i++
		if i < len(content) && content[i] == ' ' {
			i++
		}
		textStart := i
		depth := 1
		for i < len(content) && depth > 0 {
			switch content[i] {
			case '[':
				depth++
			case ']':
				depth--
			}
			if depth == 0 {
				break
			}
			i++
		}
		if depth != 0 {
			return "", "", false, 0, false
		}
		return label, content[textStart:i], true, i + 1, true
	}

	return "", "", false, 0, false
}

func isLabelChar(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9' || b == '_' || b == '-'
}
