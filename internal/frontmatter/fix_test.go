package frontmatter

import "testing"

func TestFixConvertsYAMLPageToTOML(t *testing.T) {
	in := "---\n" +
		"date: 2026-07-16T23:58:23-04:00\n" +
		"lastmod: 2026-07-16T23:58:23-04:00\n" +
		"title: Friends\n" +
		"draft: false\n" +
		"slug: friends\n" +
		"---\n\n" +
		"Here's some of my friends...\n"

	res, err := Fix("content/pages/friends.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	want := "+++\n" +
		"date = 2026-07-16T23:58:23-04:00\n" +
		"lastmod = 2026-07-16T23:58:23-04:00\n" +
		"title = 'Friends'\n" +
		"draft = false\n" +
		"slug = 'friends'\n" +
		"+++\n\n" +
		"Here's some of my friends...\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
	if !res.Changed {
		t.Fatal("expected Changed = true")
	}
}

func TestFixConvertsScriptureHeavyYAMLPostToTOML(t *testing.T) {
	// The real risk case: a post with a mixed-quote scripture summary, a
	// tags array, and a long query-string cover URL.
	in := "---\n" +
		"date: 2025-12-28T07:31:03.000Z\n" +
		"lastmod: 2026-01-06T07:29:43.000Z\n" +
		"title: Apophasis\n" +
		"draft: false\n" +
		"slug: apophasis\n" +
		`t: ["Writings","Spiritual Life"]` + "\n" +
		"cover: https://images.unsplash.com/photo-1716?crop=entropy&cs=tinysrgb&w=2000\n" +
		`cover-alt: "Photo by [Coralie](https://unsplash.com/@zenphotos) / [Unsplash](https://unsplash.com/)"` + "\n" +
		`summary: '"For the moment all discipline seems painful rather than pleasant." - Hebrews 12:11 (RSV)'` + "\n" +
		"---\n\nbody text\n"

	res, err := Fix("content/posts/2025/12/28/apophasis/index.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	want := "+++\n" +
		"date = 2025-12-28T07:31:03.000Z\n" +
		"lastmod = 2026-01-06T07:29:43.000Z\n" +
		"title = 'Apophasis'\n" +
		"draft = false\n" +
		"slug = 'apophasis'\n" +
		"t = ['Writings', 'Spiritual Life']\n" +
		"cover = 'https://images.unsplash.com/photo-1716?crop=entropy&cs=tinysrgb&w=2000'\n" +
		"cover-alt = 'Photo by [Coralie](https://unsplash.com/@zenphotos) / [Unsplash](https://unsplash.com/)'\n" +
		`summary = '"For the moment all discipline seems painful rather than pleasant." - Hebrews 12:11 (RSV)'` + "\n" +
		"+++\n\nbody text\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
}

func TestFixAlreadyCanonicalTOMLIsNoOp(t *testing.T) {
	in := "+++\n" +
		"date = 2025-12-28T07:31:03.000Z\n" +
		"lastmod = 2026-01-06T07:29:43.000Z\n" +
		"title = 'Apophasis'\n" +
		"draft = false\n" +
		"slug = 'apophasis'\n" +
		"t = ['Writings', 'Spiritual Life']\n" +
		`summary = '"quote" - Hebrews'` + "\n" +
		"+++\n\nbody text\n"

	res, err := Fix("content/posts/2025/12/28/apophasis/index.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Fatalf("expected no change, got:\n%s", res.Body)
	}
	if string(res.Body) != in {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, in)
	}
}

func TestFixReordersFieldsToSchema(t *testing.T) {
	// Fields out of order and TOML already; should just reorder.
	in := "+++\n" +
		"title = 'Apophasis'\n" +
		"date = 2025-12-28T07:31:03.000Z\n" +
		"draft = false\n" +
		"lastmod = 2026-01-06T07:29:43.000Z\n" +
		"slug = 'apophasis'\n" +
		"+++\n\nbody\n"

	res, err := Fix("content/posts/2025/12/28/apophasis/index.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	want := "+++\n" +
		"date = 2025-12-28T07:31:03.000Z\n" +
		"lastmod = 2026-01-06T07:29:43.000Z\n" +
		"title = 'Apophasis'\n" +
		"draft = false\n" +
		"slug = 'apophasis'\n" +
		"+++\n\nbody\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
}

func TestFixFillsLastmodAndSlugDefaults(t *testing.T) {
	in := "+++\ndate = 2026-07-17T00:00:00.000Z\ntitle = 'A New Post'\ndraft = true\n+++\n\nbody\n"
	res, err := Fix("content/posts/2026/07/17/a-new-post/index.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	want := "+++\n" +
		"date = 2026-07-17T00:00:00.000Z\n" +
		"lastmod = 2026-07-17T00:00:00.000Z\n" +
		"title = 'A New Post'\n" +
		"draft = true\n" +
		"slug = 'a-new-post'\n" +
		"+++\n\nbody\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
}

func TestFixDoesNotFabricateSummary(t *testing.T) {
	in := "+++\ndate = 2026-07-17T00:00:00.000Z\nlastmod = 2026-07-17T00:00:00.000Z\ntitle = 'Bare'\ndraft = true\nslug = 'bare'\n+++\n\nbody\n"
	res, err := Fix("content/pages/bare.md", []byte(in), FormatTOML)
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Fatalf("expected no change (no summary invented), got:\n%s", res.Body)
	}
}

func TestFixToYAMLTarget(t *testing.T) {
	in := "+++\ndate = 2026-07-16T23:58:23-04:00\nlastmod = 2026-07-16T23:58:23-04:00\ntitle = 'Friends'\ndraft = false\nslug = 'friends'\n+++\n\nbody\n"
	res, err := Fix("content/pages/friends.md", []byte(in), FormatYAML)
	if err != nil {
		t.Fatal(err)
	}
	want := "---\n" +
		"date: 2026-07-16T23:58:23-04:00\n" +
		"lastmod: 2026-07-16T23:58:23-04:00\n" +
		"title: Friends\n" +
		"draft: false\n" +
		"slug: friends\n" +
		"---\n\nbody\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
}

func TestFixNoFrontMatterErrors(t *testing.T) {
	_, err := Fix("content/pages/x.md", []byte("just body text\n"), FormatTOML)
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestFixUnparseableLineErrors(t *testing.T) {
	in := "---\ntitle Home\n---\n\nbody\n"
	_, err := Fix("content/pages/x.md", []byte(in), FormatTOML)
	if err == nil {
		t.Fatal("expected an error")
	}
}
