package footnotes

import "testing"

func TestRenumberInlineShorthand(t *testing.T) {
	in := "First point.[^meow: A cat sound.] Second point.[^woof: A dog sound.]\n" +
		"Referenced again here.[^meow]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Missing) != 0 {
		t.Fatalf("unexpected missing: %v", res.Missing)
	}
	want := "First point.[^1] Second point.[^2]\n" +
		"Referenced again here.[^1]\n" +
		"\n" +
		"[^1]: A cat sound.\n" +
		"[^2]: A dog sound.\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
	if !res.Changed {
		t.Fatal("expected Changed = true")
	}
}

func TestRenumberBareRefWithTrailingDef(t *testing.T) {
	in := "Some text.[^note]\n\n[^note]: The definition.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "Some text.[^1]\n\n[^1]: The definition.\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}

func TestRenumberIdempotentOnAlreadyClean(t *testing.T) {
	in := "One.[^1] Two.[^2]\n\n[^1]: First.\n[^2]: Second.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed {
		t.Fatalf("expected no change, got:\n%q", res.Body)
	}
	if res.Body != in {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, in)
	}
}

func TestRenumberMissingDefinition(t *testing.T) {
	in := "Some text.[^orphan]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "orphan" {
		t.Fatalf("expected [orphan] missing, got %v", res.Missing)
	}
	if res.Changed {
		t.Fatal("expected no change when a definition is missing")
	}
}

func TestRenumberLinkInsideInlineShorthand(t *testing.T) {
	in := "See this.[^cite: cf. [a link](https://example.com/page) for more.]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Missing) != 0 {
		t.Fatalf("unexpected missing: %v", res.Missing)
	}
	want := "See this.[^1]\n\n[^1]: cf. [a link](https://example.com/page) for more.\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}

func TestRenumberReordersByFirstAppearance(t *testing.T) {
	in := "B first.[^b: bee] A second.[^a: aye]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "B first.[^1] A second.[^2]\n\n[^1]: bee\n[^2]: aye\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}

func TestRenumberNoFootnotesUnchanged(t *testing.T) {
	in := "Just plain text, no notes at all.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if res.Changed || res.Body != in {
		t.Fatalf("expected untouched, got:\n%q", res.Body)
	}
}
