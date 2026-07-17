package frontmatter

import "testing"

func TestFixConvertsTOMLPageToCanonicalYAML(t *testing.T) {
	in := "+++\n" +
		"date = '2026-07-16T23:58:23-04:00'\n" +
		"draft = false\n" +
		"title = 'Friends'\n" +
		"slug = 'friends'\n" +
		"+++\n\n" +
		"Here's some of my friends...\n"

	res, err := Fix("content/pages/friends.md", []byte(in))
	if err != nil {
		t.Fatal(err)
	}
	want := "---\n" +
		"date: 2026-07-16T23:58:23-04:00\n" +
		"lastmod: 2026-07-16T23:58:23-04:00\n" +
		"title: Friends\n" +
		"draft: false\n" +
		"slug: friends\n" +
		"---\n\n" +
		"Here's some of my friends...\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
	if !res.Changed {
		t.Fatal("expected Changed = true")
	}
}

func TestFixAlreadyCleanPageIsNoOp(t *testing.T) {
	in := "---\n" +
		"date: 2025-03-05T20:49:24.000Z\n" +
		"lastmod: 2026-06-14T20:12:46.000Z\n" +
		"title: About\n" +
		"draft: false\n" +
		"slug: about\n" +
		"---\n\n" +
		"I'm Rdr. Seraphim Pardee...\n"

	res, err := Fix("content/pages/about.md", []byte(in))
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

func TestFixAlreadyCleanPostIsNoOp(t *testing.T) {
	in := "---\n" +
		"date: 2025-12-28T07:31:03.000Z\n" +
		"lastmod: 2026-01-06T07:29:43.000Z\n" +
		"title: Apophasis\n" +
		"draft: false\n" +
		"slug: apophasis\n" +
		`t: ["Writings","Spiritual Life"]` + "\n" +
		"cover: https://example.com/cover.jpg\n" +
		`cover-alt: "Photo by [Coralie](https://unsplash.com/@zenphotos)"` + "\n" +
		`summary: '"quote" - Hebrews 12:11 (RSV)'` + "\n" +
		"---\n\n" +
		"body text\n"

	res, err := Fix("content/posts/2025/12/28/apophasis/index.md", []byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Fatalf("expected no change, got:\n%s", res.Body)
	}
}

func TestFixPostMissingCoverAndTagsKeepsThemAbsent(t *testing.T) {
	in := "---\n" +
		"date: 1970-01-01T00:00:00.000Z\n" +
		"lastmod: 2025-06-04T15:56:58.000Z\n" +
		"title: Generative AI and Spirituality\n" +
		"draft: true\n" +
		"slug: generative-ai-and-spirituality\n" +
		"cover: /image.png\n" +
		"summary: 'A summary.'\n" +
		"---\n\nbody\n"

	res, err := Fix("content/posts/drafts/generative-ai-and-spirituality/index.md", []byte(in))
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Fatalf("expected no change, got:\n%s", res.Body)
	}
}

func TestFixDerivesSlugFromPageBundleDirectory(t *testing.T) {
	in := "---\ndate: 2026-07-17T00:00:00.000Z\ntitle: A New Post\ndraft: true\n---\n\nbody\n"
	res, err := Fix("content/posts/2026/07/17/a-new-post/index.md", []byte(in))
	if err != nil {
		t.Fatal(err)
	}
	want := "---\n" +
		"date: 2026-07-17T00:00:00.000Z\n" +
		"lastmod: 2026-07-17T00:00:00.000Z\n" +
		"title: A New Post\n" +
		"draft: true\n" +
		"slug: a-new-post\n" +
		"---\n\nbody\n"
	if string(res.Body) != want {
		t.Fatalf("got:\n%s\nwant:\n%s", res.Body, want)
	}
}

func TestFixUnknownSectionOnlyFillsDefaults(t *testing.T) {
	in := "---\ntitle: Home\ntype: home\ncover: /images/me.jpg\n---\n\n"
	res, err := Fix("content/_index.md", []byte(in))
	if err != nil {
		t.Fatal(err)
	}
	// no "date" present, so lastmod can't be defaulted; _index.md is a
	// section/home index, which conventionally doesn't carry its own
	// slug, so none is added either. Nothing to do.
	if res.Changed {
		t.Fatalf("expected no change, got:\n%s", res.Body)
	}
}

func TestFixNoFrontMatterErrors(t *testing.T) {
	_, err := Fix("content/pages/x.md", []byte("just body text\n"))
	if err == nil {
		t.Fatal("expected an error")
	}
}

func TestFixUnparseableLineErrors(t *testing.T) {
	in := "---\ntitle Home\n---\n\nbody\n"
	_, err := Fix("content/pages/x.md", []byte(in))
	if err == nil {
		t.Fatal("expected an error")
	}
}
