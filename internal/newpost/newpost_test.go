package newpost

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSlugify(t *testing.T) {
	cases := map[string]string{
		"No, You Don't Need a Significant Other": "no-you-dont-need-a-significant-other",
		"Apophasis":                              "apophasis",
		"  leading and trailing  ":               "leading-and-trailing",
	}
	for in, want := range cases {
		if got := Slugify(in); got != want {
			t.Errorf("Slugify(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFrontMatterMatchesExistingSchema(t *testing.T) {
	date := time.Date(2025, 12, 28, 7, 31, 3, 0, time.UTC)
	p := Post{
		Title:    "Apophasis",
		Slug:     "apophasis",
		Tags:     []string{"Writings", "Spiritual Life"},
		Summary:  `"quote" - Hebrews`,
		Cover:    "https://example.com/cover.jpg",
		CoverAlt: "Photo by someone",
		Draft:    false,
		Date:     date,
	}
	got := FrontMatter(p)
	want := "date: 2025-12-28T07:31:03.000Z\n" +
		"lastmod: 2025-12-28T07:31:03.000Z\n" +
		"title: Apophasis\n" +
		"draft: false\n" +
		"slug: apophasis\n" +
		"t: [\"Writings\",\"Spiritual Life\"]\n" +
		"cover: https://example.com/cover.jpg\n" +
		"cover-alt: \"Photo by someone\"\n" +
		"summary: \"\\\"quote\\\" - Hebrews\"\n"
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestFrontMatterOmitsEmptyOptionalFields(t *testing.T) {
	p := Post{Title: "Bare", Slug: "bare", Draft: true, Date: time.Now()}
	got := FrontMatter(p)
	for _, key := range []string{"\nt:", "\ncover:", "\ncover-alt:", "\nsummary:"} {
		if strings.Contains(got, key) {
			t.Errorf("expected no %q line, got:\n%s", key, got)
		}
	}
}

func TestCreateWritesPageBundleAndRejectsOverwrite(t *testing.T) {
	dir := t.TempDir()
	postsDir := filepath.Join(dir, "content", "posts")

	p := Post{
		Title: "My New Post",
		Slug:  "my-new-post",
		Draft: true,
		Date:  time.Date(2026, 7, 17, 12, 0, 0, 0, time.UTC),
	}

	path, err := Create(postsDir, p)
	if err != nil {
		t.Fatal(err)
	}
	wantPath := filepath.Join(postsDir, "2026", "07", "17", "my-new-post", "index.md")
	if path != wantPath {
		t.Fatalf("path = %q, want %q", path, wantPath)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(data), "---\n") || !strings.Contains(string(data), "title: My New Post\n") {
		t.Fatalf("unexpected content:\n%s", data)
	}

	if _, err := Create(postsDir, p); err == nil {
		t.Fatal("expected error when creating a post that already exists")
	}
}
