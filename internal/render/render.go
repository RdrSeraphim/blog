// Package render turns a post's Markdown body into standalone HTML for
// terminal-external preview, roughly mirroring the goldmark extensions
// enabled in hugo.toml (GFM tables/autolinks, footnotes, unsafe raw HTML).
package render

import (
	"bytes"
	"fmt"
	"html"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	goldmarkhtml "github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM, extension.Footnote),
	goldmark.WithParserOptions(parser.WithAutoHeadingID()),
	goldmark.WithRendererOptions(goldmarkhtml.WithUnsafe()),
)

// ToHTML converts a Markdown body to an HTML fragment.
func ToHTML(body string) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert([]byte(body), &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Page wraps an HTML fragment in a minimal, theme-aware standalone page.
func Page(title, fragment string) string {
	if title == "" {
		title = "Preview"
	}
	return fmt.Sprintf(pageTemplate, html.EscapeString(title), fragment)
}

const pageTemplate = `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>%s</title>
<style>
  :root { color-scheme: light dark; }
  body {
    max-width: 42rem;
    margin: 3rem auto;
    padding: 0 1.25rem 4rem;
    font: 1.05rem/1.7 -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
    color: #1a1a1a;
    background: #fdfdfb;
  }
  @media (prefers-color-scheme: dark) {
    body { color: #e6e6e6; background: #14161a; }
    a { color: #7db4ff; }
    blockquote { border-left-color: #444; color: #aaa; }
    code, pre { background: #1e2126; }
    hr { border-color: #333; }
  }
  h1, h2, h3 { line-height: 1.25; }
  a { color: #0552b5; }
  blockquote { margin: 1.25rem 0; padding: 0 1rem; border-left: 3px solid #ccc; color: #555; }
  code { padding: .1em .3em; border-radius: 4px; background: #f0f0f0; font-size: .9em; }
  pre { padding: 1rem; overflow-x: auto; border-radius: 6px; background: #f0f0f0; }
  pre code { padding: 0; background: none; }
  img { max-width: 100%%; }
  hr { border: none; border-top: 1px solid #ddd; margin: 2rem 0; }
  .footnotes { font-size: .9em; color: #666; }
  table { border-collapse: collapse; }
  td, th { border: 1px solid #ccc; padding: .3em .6em; }
</style>
</head>
<body>
%s
</body>
</html>
`
