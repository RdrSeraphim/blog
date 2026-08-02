# blog

The musings of an Orthodox Christian reader and software engineer.

## Authoring

- **New posts/pages**: `hugo new posts/YYYY/MM/DD/slug/index.md` or `hugo new pages/slug.md`. The archetypes in `archetypes/` (`posts.md`, `pages.md`, and a `default.md` fallback for any other section) fill in TOML front matter matching the fields used across existing content — title-cased title, slug, date/lastmod, `draft = true`, and empty `summary` (plus `t` for posts) placeholders.
- **Preview**: `npm run dev` (`hugo server --enableGitInfo --port 8787 --buildDrafts --disableFastRender`) for a full, live-reloading themed preview in the browser. For in-terminal rendering, use an editor plugin like [render-markdown.nvim](https://github.com/MeanderingProgrammer/render-markdown.nvim)'s `:RenderMarkdown`.
- **Scaffolding, footnotes, front matter, and linting**: see [`blogcli`](https://github.com/RdrSeraphim/blogcli), a separate tool — its README has usage.

Front matter across the site is **TOML (`+++`)** — matching `hugo.toml` and Hugo's own defaults. (The bulk of the content originally arrived as YAML from a Ghost migration; `blogcli frontmatter` converted it.)

## Content linting

Front matter linting (via `blogcli lint`) runs in CI (`.github/workflows/lint-content.yml`) and is available as an opt-in pre-commit hook:

```sh
git config core.hooksPath .githooks
```

See [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for Neovim keymaps wired to `blogcli`'s commands.
