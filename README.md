# blog

The musings of an Orthodox Christian reader and software engineer.

## Authoring

- **New posts/pages**: `hugo new posts/YYYY/MM/DD/slug/index.md` or `hugo new pages/slug.md`. The archetypes in `archetypes/` (`posts.md`, `pages.md`, and a `default.md` fallback for any other section) fill in TOML front matter matching the fields used across existing content — title-cased title, slug, date/lastmod, `draft = true`, and empty `summary` (plus `t` for posts) placeholders.
- **Preview**: `npm run dev` (`hugo server --enableGitInfo --port 8787 --buildDrafts --disableFastRender`) for a full, live-reloading themed preview in the browser. For in-terminal rendering, use an editor plugin like [render-markdown.nvim](https://github.com/MeanderingProgrammer/render-markdown.nvim)'s `:RenderMarkdown`.
- **Footnotes, front matter, and linting**: see `cmd/blogcli` below.

Front matter across the site is **TOML (`+++`)** — matching `hugo.toml` and Hugo's own defaults. (The bulk of the content originally arrived as YAML from a Ghost migration; `blogcli frontmatter` converted it.)

## blogcli

`cmd/blogcli` handles the authoring chores Hugo and editor plugins don't cover:

```sh
go build -o bin/blogcli ./cmd/blogcli
bin/blogcli footnotes   [-w] <file>
bin/blogcli frontmatter [-w] [--yaml] [<file|dir>...]
bin/blogcli lint             [--strict] [<file|dir>...]
```

`footnotes` and `frontmatter` print to stdout by default; pass `-w` to write back. With no path arguments, `frontmatter` and `lint` process every `.md` under `./content` (a directory is walked recursively).

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

`blogcli frontmatter -w [path...]` normalizes content front matter to the site's conventions:

- re-serializes to canonical **TOML** (`--yaml` goes the other way, if the canonical format ever changes)
- fills in a missing `lastmod` from `date`, and a missing `slug` from the file's path
- for files under `content/posts/` or `content/pages/`, reorders fields to match the schema used across the rest of that section

It never invents or drops a value beyond those two safe defaults — in particular it never fabricates a `summary` — and leaves fields it doesn't recognize as-is. Already-canonical files are a no-op.

### lint

`blogcli lint [path...]` checks front matter without changing anything and reports issues per file. **Errors** (non-canonical format, unparseable front matter, a missing required field, an invalid date) exit non-zero; **warnings** (a missing/empty `summary`, so the meta description is uncurated; a `cover` without a `cover-alt`; a `slug` that doesn't match the file's path) are advisory unless you pass `--strict`.

This runs in CI (`.github/workflows/lint-content.yml`) and is available as an opt-in pre-commit hook:

```sh
git config core.hooksPath .githooks
```

See [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for Neovim keymaps wired to these commands.
