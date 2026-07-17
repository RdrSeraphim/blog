package footnotes

import "testing"

func TestRenumberContentIsTheNote(t *testing.T) {
	in := "First point.[^A cat sound.] Second point.[^A dog sound.]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Missing) != 0 {
		t.Fatalf("unexpected missing: %v", res.Missing)
	}
	want := "First point.[^1] Second point.[^2]\n" +
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

func TestRenumberRepeatedTextIsNotDeduped(t *testing.T) {
	// Two markers with identical text are still two distinct footnotes.
	in := "One.[^same] Two.[^same]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "One.[^1] Two.[^2]\n\n[^1]: same\n[^2]: same\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}

func TestRenumberLegacyNumericRefWithTrailingDef(t *testing.T) {
	in := "Some text.[^7]\n\n[^7]: The definition.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "Some text.[^1]\n\n[^1]: The definition.\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}

func TestRenumberLegacyNumericRefRepeatedIsDeduped(t *testing.T) {
	// Multiple markers pointing at the same original number are the same footnote.
	in := "One.[^1] Again.[^1]\n\n[^1]: Shared definition.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "One.[^1] Again.[^1]\n\n[^1]: Shared definition.\n"
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

func TestRenumberMissingLegacyDefinition(t *testing.T) {
	in := "Some text.[^9]\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Missing) != 1 || res.Missing[0] != "9" {
		t.Fatalf("expected [9] missing, got %v", res.Missing)
	}
	if res.Changed {
		t.Fatal("expected no change when a definition is missing")
	}
}

func TestRenumberLinkInsideNoteText(t *testing.T) {
	in := "See this.[^cf. [a link](https://example.com/page) for more.]\n"
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
	in := "B first.[^bee] A second.[^aye]\n"
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

func TestRenumberMixedLegacyAndLiteral(t *testing.T) {
	// A post being edited that already has a numbered footnote at the
	// bottom, plus a brand new one just written inline.
	in := "Old point.[^1] New point.[^Just added this.]\n\n[^1]: Original text.\n"
	res, err := Renumber(in)
	if err != nil {
		t.Fatal(err)
	}
	want := "Old point.[^1] New point.[^2]\n\n[^1]: Original text.\n[^2]: Just added this.\n"
	if res.Body != want {
		t.Fatalf("got:\n%q\nwant:\n%q", res.Body, want)
	}
}
