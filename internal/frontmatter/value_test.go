package frontmatter

import "testing"

func TestDecodeEncodeRoundTripsTOML(t *testing.T) {
	// Every value, once in canonical TOML form, must survive a
	// decode->encode round trip unchanged (this is what makes running
	// `frontmatter` on an already-canonical file a no-op).
	cases := []string{
		`false`,
		`true`,
		`2025-12-28T07:31:03.000Z`,
		`2026-07-16T23:58:23-04:00`,
		`'Friends'`,
		`'A Grief Unobserved'`,
		`"Y'ever been stuck with some monks in a random U-Haul center in Ohio?"`,
		`'"For the moment all discipline seems painful." - Hebrews 12:11 (RSV)'`,
		`'https://images.unsplash.com/photo-1716?crop=entropy&cs=tinysrgb&w=2000'`,
		`[]`,
		`['Writings']`,
		`['Writings', 'Spiritual Life']`,
		`''`,
	}
	for _, in := range cases {
		got := DecodeValue(in, FormatTOML).Encode(FormatTOML)
		if got != in {
			t.Errorf("round trip changed %q -> %q", in, got)
		}
	}
}

func TestYAMLScriptureSummaryToTOMLKeepsDoubleQuotes(t *testing.T) {
	// Single-quoted YAML wrapping a scripture quote (double quotes inside,
	// no single quotes) should become a TOML literal string, preserving the
	// inner double quotes verbatim.
	in := `'"Better to go to the house of mourning than to go to the house of feasting." - Ecclesiastes 7:2'`
	v := DecodeValue(in, FormatYAML)
	got := v.Encode(FormatTOML)
	want := `'"Better to go to the house of mourning than to go to the house of feasting." - Ecclesiastes 7:2'`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestYAMLApostropheSummaryToTOMLUsesBasicString(t *testing.T) {
	// Double-quoted YAML wrapping text with an apostrophe (single quote,
	// no double quotes) must become a TOML basic string, since a literal
	// string can't hold the apostrophe.
	in := `"Y'ever been stuck with some monks in a random U-Haul center in Ohio?"`
	got := DecodeValue(in, FormatYAML).Encode(FormatTOML)
	want := `"Y'ever been stuck with some monks in a random U-Haul center in Ohio?"`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestValueWithBothQuoteKindsToTOML(t *testing.T) {
	// A value containing BOTH a double quote and an apostrophe can't be a
	// literal string; it must be a basic string with the double quote
	// escaped and the apostrophe left as-is.
	v := Value{kind: kindString, str: `He said "hi" and it's fine`}
	got := v.Encode(FormatTOML)
	want := `"He said \"hi\" and it's fine"`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
	// ...and it decodes back to the original text.
	back := DecodeValue(got, FormatTOML)
	if back.str != v.str {
		t.Fatalf("decode of %q gave %q, want %q", got, back.str, v.str)
	}
}

func TestTitleWithColonAndAmpersandToTOML(t *testing.T) {
	in := `"Reflection: Age & Legacy (pt. 1)"`
	got := DecodeValue(in, FormatYAML).Encode(FormatTOML)
	want := `'Reflection: Age & Legacy (pt. 1)'`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestUnquotedYAMLTitleWithApostropheToTOML(t *testing.T) {
	// "No, You Don't Need a Significant Other" is stored unquoted in YAML.
	in := `No, You Don't Need a Significant Other`
	got := DecodeValue(in, FormatYAML).Encode(FormatTOML)
	want := `"No, You Don't Need a Significant Other"`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestDateTimeTokenPreservedExactly(t *testing.T) {
	in := `2025-12-28T07:31:03.000Z`
	v := DecodeValue(in, FormatYAML)
	if v.kind != kindDateTime {
		t.Fatalf("expected datetime, got kind %d", v.kind)
	}
	if got := v.Encode(FormatTOML); got != in {
		t.Fatalf("datetime token changed: %q -> %q", in, got)
	}
}

func TestArrayConversionYAMLToTOMLAddsSpaces(t *testing.T) {
	in := `["Writings","Spiritual Life"]`
	got := DecodeValue(in, FormatYAML).Encode(FormatTOML)
	want := `['Writings', 'Spiritual Life']`
	if got != want {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}
