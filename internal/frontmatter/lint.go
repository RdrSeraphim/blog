package frontmatter

import (
	"fmt"
	"sort"
	"strings"
)

// Severity distinguishes issues that should fail a check run from advisory
// ones.
type Severity int

const (
	// Warning is advisory: worth fixing, but doesn't fail a check.
	Warning Severity = iota
	// Error should fail a --check run (e.g. in CI or a pre-commit hook).
	Error
)

func (s Severity) String() string {
	if s == Error {
		return "error"
	}
	return "warning"
}

// Issue is a single problem found in a content file's front matter.
type Issue struct {
	Severity Severity
	Message  string
}

// Lint inspects a content file's front matter against this blog's
// conventions and returns any issues found, most severe first. The target
// format is the canonical one the blog is standardizing on; a file in the
// other format is reported as an error (it should be converted with
// `blogcli frontmatter -w`).
func Lint(path string, data []byte, target Format) []Issue {
	var issues []Issue
	add := func(sev Severity, format string, args ...any) {
		issues = append(issues, Issue{Severity: sev, Message: fmt.Sprintf(format, args...)})
	}

	raw, _, srcTOML, ok := SplitAny(data)
	if !ok {
		add(Error, "no front matter found")
		return issues
	}
	src := FormatYAML
	if srcTOML {
		src = FormatTOML
	}
	if src != target {
		add(Error, "front matter is %s; canonical is %s (run: blogcli frontmatter -w %s)", formatName(src), formatName(target), path)
	}

	fields, err := ParseFields(raw)
	if err != nil {
		add(Error, "%v", err)
		return issues
	}

	byKey := make(map[string]Value, len(fields))
	for _, f := range fields {
		byKey[f.Key] = DecodeValue(f.Value, src)
	}

	order := sectionFieldOrder(path)
	if order == nil {
		// Not a section with an established schema; format check above is
		// all we can meaningfully say.
		return issues
	}

	for _, key := range order {
		if _, ok := byKey[key]; ok {
			continue
		}
		switch key {
		case "t", "cover", "cover-alt":
			// genuinely optional; absence is fine
		case "summary":
			add(Warning, "no summary - the theme has nothing curated to use for the meta description")
		default:
			add(Error, "missing required field %q", key)
		}
	}

	if summary, ok := byKey["summary"]; ok && strings.TrimSpace(summary.str) == "" {
		add(Warning, "summary is empty - the theme has nothing curated to use for the meta description")
	}
	if _, hasCover := byKey["cover"]; hasCover {
		if alt, ok := byKey["cover-alt"]; !ok || strings.TrimSpace(alt.str) == "" {
			add(Warning, "cover set without a cover-alt")
		}
	}
	if slug, ok := byKey["slug"]; ok {
		if want := inferSlug(path); want != "" && slug.str != want {
			add(Warning, "slug %q doesn't match the path-derived %q", slug.str, want)
		}
	}
	for _, key := range []string{"date", "lastmod"} {
		v, ok := byKey[key]
		if !ok {
			continue
		}
		if v.kind == kindDateTime || (v.kind == kindString && looksLikeDateTime(v.str)) {
			continue
		}
		add(Error, "%s is not a valid datetime: %s", key, v.str)
	}

	sort.SliceStable(issues, func(i, j int) bool {
		return issues[i].Severity > issues[j].Severity
	})
	return issues
}

func formatName(f Format) string {
	if f == FormatTOML {
		return "TOML"
	}
	return "YAML"
}
