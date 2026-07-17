package frontmatter

import (
	"strings"
	"testing"
)

func hasError(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == Error {
			return true
		}
	}
	return false
}

func TestLintCleanTOMLPostNoIssues(t *testing.T) {
	in := "+++\n" +
		"date = 2025-12-28T07:31:03.000Z\n" +
		"lastmod = 2026-01-06T07:29:43.000Z\n" +
		"title = 'Apophasis'\n" +
		"draft = false\n" +
		"slug = 'apophasis'\n" +
		"t = ['Writings']\n" +
		"cover = 'https://example.com/c.jpg'\n" +
		"cover-alt = 'Photo by someone'\n" +
		`summary = '"quote" - Hebrews'` + "\n" +
		"+++\n\nbody\n"
	issues := Lint("content/posts/2025/12/28/apophasis/index.md", []byte(in), FormatTOML)
	if len(issues) != 0 {
		t.Fatalf("expected no issues, got %v", issues)
	}
}

func TestLintFlagsYAMLWhenCanonicalIsTOML(t *testing.T) {
	in := "---\ndate: 2025-12-28T07:31:03.000Z\nlastmod: 2025-12-28T07:31:03.000Z\ntitle: X\ndraft: false\nslug: x\nsummary: 'y'\n---\n\nbody\n"
	issues := Lint("content/pages/x.md", []byte(in), FormatTOML)
	if !hasError(issues) {
		t.Fatalf("expected an error for wrong format, got %v", issues)
	}
}

func TestLintWarnsOnMissingSummary(t *testing.T) {
	in := "+++\ndate = 2025-12-28T07:31:03.000Z\nlastmod = 2025-12-28T07:31:03.000Z\ntitle = 'Friends'\ndraft = false\nslug = 'friends'\n+++\n\nbody\n"
	issues := Lint("content/pages/friends.md", []byte(in), FormatTOML)
	if hasError(issues) {
		t.Fatalf("did not expect an error, got %v", issues)
	}
	found := false
	for _, i := range issues {
		if i.Severity == Warning && strings.Contains(i.Message, "summary") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a summary warning, got %v", issues)
	}
}

func TestLintWarnsOnCoverWithoutAlt(t *testing.T) {
	in := "+++\ndate = 2025-12-28T07:31:03.000Z\nlastmod = 2025-12-28T07:31:03.000Z\ntitle = 'X'\ndraft = false\nslug = 'x'\nsummary = 's'\ncover = 'https://example.com/c.jpg'\n+++\n\nbody\n"
	issues := Lint("content/posts/2025/12/28/x/index.md", []byte(in), FormatTOML)
	found := false
	for _, i := range issues {
		if strings.Contains(i.Message, "cover-alt") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a cover-alt warning, got %v", issues)
	}
}

func TestLintWarnsOnSlugPathMismatch(t *testing.T) {
	in := "+++\ndate = 2025-12-28T07:31:03.000Z\nlastmod = 2025-12-28T07:31:03.000Z\ntitle = 'X'\ndraft = false\nslug = 'wrong'\nsummary = 's'\n+++\n\nbody\n"
	issues := Lint("content/pages/friends.md", []byte(in), FormatTOML)
	found := false
	for _, i := range issues {
		if strings.Contains(i.Message, "slug") {
			found = true
		}
	}
	if !found {
		t.Fatalf("expected a slug mismatch warning, got %v", issues)
	}
}

func TestLintErrorsOnMissingRequiredField(t *testing.T) {
	in := "+++\ndate = 2025-12-28T07:31:03.000Z\nlastmod = 2025-12-28T07:31:03.000Z\ndraft = false\nslug = 'x'\nsummary = 's'\n+++\n\nbody\n"
	issues := Lint("content/pages/x.md", []byte(in), FormatTOML)
	if !hasError(issues) {
		t.Fatalf("expected an error for missing title, got %v", issues)
	}
}
