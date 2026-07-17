# blog

The musings of an Orthodox Christian reader and software engineer.

## blogcli

`cmd/blogcli` is a small Go CLI/TUI for authoring posts:

```sh
go build -o bin/blogcli ./cmd/blogcli

bin/blogcli new                            # interactive form -> scaffolds content/posts/YYYY/MM/DD/slug/index.md
bin/blogcli footnotes [-w] <file>          # renumber mnemonic footnotes into [^1], [^2]...
bin/blogcli preview [--html] <file>        # preview in a terminal viewer, or a live-reloading browser page
```

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

See `bin/blogcli <command> -h` for details, and [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for Neovim keymaps wired to this tool.

