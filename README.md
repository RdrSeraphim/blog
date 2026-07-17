# blog

The musings of an Orthodox Christian reader and software engineer.

## Authoring

- **New posts/pages**: `hugo new posts/YYYY/MM/DD/slug/index.md` or `hugo new pages/slug.md`. The archetypes in `archetypes/` (`posts.md`, `pages.md`, and a `default.md` fallback for any other section) fill in front matter matching the fields already used across existing content — title-cased title, slug, date/lastmod, `draft: true`, and (for posts) empty `t`/`summary` placeholders.
- **Preview**: `npm run dev` (`hugo server --enableGitInfo --port 8787 --buildDrafts --disableFastRender`) for a full, live-reloading themed preview in the browser. For in-terminal rendering, use an editor plugin like [render-markdown.nvim](https://github.com/MeanderingProgrammer/render-markdown.nvim)'s `:RenderMarkdown`.
- **Footnotes**: see `cmd/blogcli` below.

## blogcli

`cmd/blogcli` renumbers hand-authored footnotes — the one piece of drafting that Hugo and editor plugins don't already cover.

```sh
go build -o bin/blogcli ./cmd/blogcli
bin/blogcli [-w] <file>
```

While drafting, write a footnote's text directly in the marker instead of tracking numbers by hand:

```markdown
Point one.[^A cat sound.] Point two.[^A dog sound.]
```

`blogcli -w <file>` rewrites that to the numbered/endnote style used across this blog:

```markdown
Point one.[^1] Point two.[^2]

[^1]: A cat sound.
[^2]: A dog sound.
```

Each marker becomes its own footnote, even if the text repeats. The one exception is a purely numeric marker like `[^1]` paired with a `[^1]: text` line elsewhere in the file — that's read as an already-numbered reference rather than literal text, so running this on an already-numbered post is a no-op.

Without `-w`, the result is printed to stdout and the file is left untouched. See [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for a Neovim keymap wired to this.
