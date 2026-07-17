// Package newpost creates new Hugo post page bundles under content/posts,
// matching the front matter schema already used across this blog.
package newpost

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// Post holds the fields collected from the user for a new post.
type Post struct {
	Title    string
	Slug     string
	Tags     []string
	Summary  string
	Cover    string
	CoverAlt string
	Draft    bool
	Date     time.Time
}

var slugInvalid = regexp.MustCompile(`[^a-z0-9]+`)

// Slugify derives a URL-safe slug from a title, e.g. "No, You Don't Need"
// becomes "no-you-dont-need".
func Slugify(title string) string {
	s := strings.ToLower(title)
	s = strings.ReplaceAll(s, "'", "")
	s = slugInvalid.ReplaceAllString(s, "-")
	return strings.Trim(s, "-")
}

// Dir returns the content/posts/YYYY/MM/DD/slug directory for a post,
// relative to the blog's contentPostsDir (e.g. "content/posts").
func Dir(contentPostsDir string, p Post) string {
	return filepath.Join(contentPostsDir,
		fmt.Sprintf("%04d", p.Date.Year()),
		fmt.Sprintf("%02d", p.Date.Month()),
		fmt.Sprintf("%02d", p.Date.Day()),
		p.Slug,
	)
}

// FrontMatter renders the YAML front matter block (without the --- delimiters).
func FrontMatter(p Post) string {
	ts := p.Date.Format("2006-01-02T15:04:05.000Z")
	var b strings.Builder
	fmt.Fprintf(&b, "date: %s\n", ts)
	fmt.Fprintf(&b, "lastmod: %s\n", ts)
	fmt.Fprintf(&b, "title: %s\n", p.Title)
	fmt.Fprintf(&b, "draft: %t\n", p.Draft)
	fmt.Fprintf(&b, "slug: %s\n", p.Slug)
	if len(p.Tags) > 0 {
		quoted := make([]string, len(p.Tags))
		for i, t := range p.Tags {
			quoted[i] = fmt.Sprintf("%q", t)
		}
		fmt.Fprintf(&b, "t: [%s]\n", strings.Join(quoted, ","))
	}
	if p.Cover != "" {
		fmt.Fprintf(&b, "cover: %s\n", p.Cover)
	}
	if p.CoverAlt != "" {
		fmt.Fprintf(&b, "cover-alt: %q\n", p.CoverAlt)
	}
	if p.Summary != "" {
		fmt.Fprintf(&b, "summary: %q\n", p.Summary)
	}
	return b.String()
}

// Create writes a new index.md page bundle for p under contentPostsDir and
// returns the path to the created file. It refuses to overwrite an
// existing post directory.
func Create(contentPostsDir string, p Post) (string, error) {
	dir := Dir(contentPostsDir, p)
	if _, err := os.Stat(dir); err == nil {
		return "", fmt.Errorf("%s already exists", dir)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "index.md")
	content := "---\n" + FrontMatter(p) + "---\n\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return "", err
	}
	return path, nil
}
