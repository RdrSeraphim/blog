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

While drafting, write footnotes with an inline mnemonic label instead of tracking numbers by hand:

```markdown
Point one.[^meow: A cat sound.] Point two, mentioned again[^meow].
```

`blogcli footnotes -w <file>` rewrites that to the numbered/endnote style used across this blog:

```markdown
Point one.[^1] Point two, mentioned again[^1].

[^1]: A cat sound.
```

A bare `[^label]` paired with a `[^label]: text` line elsewhere in the file works too, so running it on an already-numbered post is a no-op.

See `bin/blogcli <command> -h` for details, and [`RdrSeraphim/kickstart.nvim`](https://github.com/rdrseraphim/kickstart.nvim) for Neovim keymaps wired to this tool.

