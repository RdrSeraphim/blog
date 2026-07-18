// Package footnotes renumbers hand-authored footnotes into the sequential
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
	Missing []string // legacy numeric references with no matching definition
}

var trailingDefRe = regexp.MustCompile(`^\[\^([^\]\s]+)\]:[ \t]?(.*)$`)
var allDigits = regexp.MustCompile(`^[0-9]+$`)

// Renumber rewrites footnotes in body to sequential numeric identifiers, in
// order of first appearance, and appends a clean endnote block containing
// their definitions.
//
// The bracket content of a footnote marker is its text, written right where
// it's used:
//
//	Point one.[^A cat sound.] Point two.[^A dog sound.]
//
// becomes:
//
//	Point one.[^1] Point two.[^2]
//
//	[^1]: A cat sound.
//	[^2]: A dog sound.
//
// Each marker stands alone: identical text in two different markers still
// produces two separate footnotes, never a shared one.
//
// The one exception is a marker whose content is purely numeric, e.g.
// [^1] paired with a "[^1]: text" line elsewhere in the document. That's
// read as an already-numbered reference rather than literal text "1", so
// re-running Renumber on an already-clean, already-published post (which
// uses exactly that numeric style) is a no-op, and multiple [^1] markers
// referring to the same definition are recognized as the same footnote.
//
// A numeric marker with no matching definition is left untouched in the
// body and reported in Missing; callers should treat a non-empty Missing
// as reason not to write the result out.
func Renumber(body string) (Result, error) {
	content, trailingDefs := extractTrailingDefs(body)

	var out strings.Builder
	order := []string{}
	numOf := map[string]int{}
	defText := map[string]string{}
	missingSet := map[string]bool{}
	occurrence := 0

	i := 0
	for i < len(content) {
		start := strings.Index(content[i:], "[^")
		if start == -1 {
			out.WriteString(content[i:])
			break
		}
		start += i
		out.WriteString(content[i:start])

		text, numeric, tokenEnd, isToken := scanToken(content, start)
		if !isToken {
			out.WriteString("[^")
			i = start + 2
			continue
		}

		var key string
		if numeric {
			// Legacy numeric reference: dedup by the original number and
			// resolve its text from the trailing definition block.
			key = "n:" + text
			if _, known := defText[key]; !known {
				if t, found := trailingDefs[text]; found {
					defText[key] = t
				} else {
					missingSet[text] = true
					out.WriteString(content[start:tokenEnd])
					i = tokenEnd
					continue
				}
			}
		} else {
			// The bracket content is the footnote's text, verbatim. Every
			// occurrence is its own footnote, even if the text repeats.
			occurrence++
			key = fmt.Sprintf("c:%d", occurrence)
			defText[key] = text
		}

		if _, seen := numOf[key]; !seen {
			numOf[key] = len(order) + 1
			order = append(order, key)
		}
		fmt.Fprintf(&out, "[^%d]", numOf[key])
		i = tokenEnd
	}

	if len(missingSet) > 0 {
		missing := make([]string, 0, len(missingSet))
		for l := range missingSet {
			missing = append(missing, l)
		}
		sort.Strings(missing)
		return Result{Body: body, Changed: false, Missing: missing}, nil
	}

	newBody := strings.TrimRight(out.String(), "\n") + "\n"
	if len(order) > 0 {
		var defs strings.Builder
		for _, key := range order {
			fmt.Fprintf(&defs, "[^%d]: %s\n", numOf[key], defText[key])
		}
		newBody += "\n" + strings.TrimRight(defs.String(), "\n") + "\n"
	}

	changed := newBody != body
	return Result{Body: newBody, Changed: changed}, nil
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
// (which must be "[^"), reading up to the balanced closing bracket so
// footnote text may itself contain markdown links. It returns the bracket
// content, whether that content is purely numeric, and the index just past
// the closing bracket. isToken is false for an empty or unbalanced "[^...]",
// in which case only the leading "[^" should be consumed by the caller.
func scanToken(content string, start int) (text string, numeric bool, end int, isToken bool) {
	i := start + 2
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
	if depth != 0 || i == textStart {
		return "", false, 0, false
	}
	text = content[textStart:i]
	return text, allDigits.MatchString(text), i + 1, true
}
