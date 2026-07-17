# blog

The musings of an Orthodox Christian reader and software engineer.

## Authoring

- **New posts/pages**: `hugo new posts/YYYY/MM/DD/slug/index.md` or `hugo new pages/slug.md`. The archetypes in `archetypes/` (`posts.md`, `pages.md`, and a `default.md` fallback for any other section) fill in front matter matching the fields already used across existing content — title-cased title, slug, date/lastmod, `draft: true`, and (for posts) empty `t`/`summary` placeholders.
- **Preview**: `npm run dev` (`hugo server --enableGitInfo --port 8787 --buildDrafts --disableFastRender`) for a full, live-reloading themed preview in the browser. For in-terminal rendering, use an editor plugin like [render-markdown.nvim](https://github.com/MeanderingProgrammer/render-markdown.nvim)'s `:RenderMarkdown`.
- **Footnotes and front matter cleanup**: see `cmd/blogcli` below.

## blogcli

`cmd/blogcli` handles the two bits of drafting that Hugo and editor plugins don't already cover:

```sh
go build -o bin/blogcli ./cmd/blogcli
bin/blogcli footnotes [-w] <file>
bin/blogcli frontmatter [-w] <file>
```

Both print the result to stdout by default; pass `-w` to write it back to the file instead.

### footnotes

While drafting, write a footnote's text directly in the marker instead of tracking numbers by hand:

```markdown
Point one.[^A cat sound.] Point two.[^A dog sound.]
```

`blogcli footnotes -w <file>` rewrites that to the numbered/endnote style used across this blog:

```markdown
Point one.[^1] Point two.[^2]

[^1]: A cat sound.
[^2]: A dog sound.
```

Each marker becomes its own footnote, even if the text repeats. The one exception is a purely numeric marker like `[^1]` paired with a `[^1]: text` line elsewhere in the file — that's read as an already-numbered reference rather than literal text, so running this on an already-numbered post is a no-op.

### frontmatter

For content that predates a proper archetype, or was hand-created with `hugo new` before `archetypes/posts.md`/`pages.md` existed (plain TOML `+++` front matter, missing `lastmod`, missing `slug`), `blogcli frontmatter -w <file>` brings it in line with the rest of the site:

- converts TOML (`+++`) front matter to this blog's YAML (`---`) style
- fills in a missing `lastmod` from `date`, and a missing `slug` from the file's path
- for files under `content/posts/` or `content/pages/`, reorders fields to match the schema used across the rest of that section

It never invents or drops a value beyond those two safe defaults, and leaves fields it doesn't recognize as they are. Already-clean files are a no-op.

See [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for Neovim keymaps wired to both commands.
