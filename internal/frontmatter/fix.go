package frontmatter

import (
	"errors"
	"path/filepath"
	"regexp"
	"strings"
)

// ErrNoFrontMatter is returned by Fix when a file has no recognizable
// front matter block at all.
var ErrNoFrontMatter = errors.New("no front matter found")

var postsFieldOrder = []string{"date", "lastmod", "title", "draft", "slug", "t", "cover", "cover-alt", "summary"}
var pagesFieldOrder = []string{"date", "lastmod", "title", "draft", "slug"}

// sectionFieldOrder returns the canonical field order for a content file,
// inferred from its path, or nil if it's not under a section this blog
// has an established schema for (in which case Fix only normalizes
// format and fills in safe defaults, without reordering fields).
func sectionFieldOrder(path string) []string {
	// Leading slash normalizes both absolute and repo-relative paths
	// (e.g. "content/pages/x.md") to the same "/content/pages/" check.
	slashed := "/" + filepath.ToSlash(path)
	switch {
	case strings.Contains(slashed, "/content/posts/"):
		return postsFieldOrder
	case strings.Contains(slashed, "/content/pages/"):
		return pagesFieldOrder
	default:
		return nil
	}
}

var slugInvalid = regexp.MustCompile(`[^a-z0-9]+`)

// inferSlug derives a slug from a content file's path the same way a slug
// would naturally come from its filename: the parent directory name for a
// leaf page bundle's index.md, otherwise the filename itself. It returns
// "" for a _index.md section/home index, which conventionally doesn't
// carry a slug of its own.
func inferSlug(path string) string {
	base := filepath.Base(path)
	switch base {
	case "_index.md":
		return ""
	case "index.md":
		base = filepath.Base(filepath.Dir(path))
	default:
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}
	s := strings.ToLower(base)
	s = strings.ReplaceAll(s, "'", "")
	s = slugInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// FixResult is the outcome of a Fix call.
type FixResult struct {
	Body    []byte // the rewritten file (unchanged if nothing to do)
	Changed bool
}

// Fix normalizes a content file's front matter: TOML (+++) blocks are
// converted to this blog's YAML style, a missing "lastmod" is filled in
// from "date", a missing "slug" is derived from the file's path, and -
// for files under content/posts or content/pages, where this blog has an
// established field schema - fields are reordered to match it. Fields
// this blog has no opinion on (an unrecognized section, or extra keys
// beyond the known schema) are left as they are, just carried over.
//
// It never invents or drops a value: only lastmod and slug get defaults,
// and only when they're missing entirely.
func Fix(path string, data []byte) (FixResult, error) {
	raw, body, _, ok := SplitAny(data)
	if !ok {
		return FixResult{}, ErrNoFrontMatter
	}
	fields, err := ParseFields(raw)
	if err != nil {
		return FixResult{}, err
	}

	byKey := make(map[string]string, len(fields))
	var original []string
	for _, f := range fields {
		byKey[f.Key] = f.Value
		original = append(original, f.Key)
	}

	for k, v := range byKey {
		byKey[k] = Unquote(v)
	}

	var added []string
	if _, ok := byKey["lastmod"]; !ok {
		if date, ok := byKey["date"]; ok {
			byKey["lastmod"] = date
			added = append(added, "lastmod")
		}
	}
	if _, ok := byKey["slug"]; !ok {
		if slug := inferSlug(path); slug != "" {
			byKey["slug"] = slug
			added = append(added, "slug")
		}
	}

	order := sectionFieldOrder(path)
	var finalOrder []string
	if order != nil {
		known := make(map[string]bool, len(order))
		for _, k := range order {
			known[k] = true
			if _, ok := byKey[k]; ok {
				finalOrder = append(finalOrder, k)
			}
		}
		for _, k := range original {
			if !known[k] {
				finalOrder = append(finalOrder, k)
			}
		}
	} else {
		finalOrder = append(finalOrder, original...)
		finalOrder = append(finalOrder, added...)
	}

	var out []Field
	for _, k := range finalOrder {
		out = append(out, Field{Key: k, Value: byKey[k]})
	}

	newBody := Join(RenderFields(out), body)
	return FixResult{Body: newBody, Changed: string(newBody) != string(data)}, nil
}
