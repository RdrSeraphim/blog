+++
date = 2025-03-05T20:49:24.000Z
lastmod = 2026-06-14T20:12:46.000Z
title = 'About'
draft = false
slug = 'about'
+++

## This Blog

This blog is a public-facing outlet for my various thoughts and tangents, including some technical postings going over software development or homelabbing.

It runs off [Hugo](https://gohugo.io), a static page generator – primarily used for blogs. Prior to this, I was using [Ghost](https://ghost.org/), but the cost and technical overhead became too great to justify keeping around. That being said, there were some things I really enjoyed out of Ghost, and so I've worked to create a largely 1-to-1 experience. The only exception being that there are no longer memberships or subscription mechanisms beyond [RSS](https://srp.life/feed.xml) and automated ActivityPub/ATProtocol accounts that link new posts when I write them ([see below](#contactsocial)).

Some things I've done:

1. The UX between the Ghost iteration and this iteration were 1-to-1, pixel perfect with path preservation. This has changed some given that I moved away from Solaris CDE styling recently, but it is still mostly the case.
2. I forked [minidon](https://github.com/yusukebe/minidon) to make [ssg-ap-bridge](https://github.com/RdrSeraphim/ssg-ap-bridge), a way for blogs powered by Static Site Generators (SSGs) to have a mechanism similar to Ghost's [ActivityPub integration](https://activitypub.ghost.org/), minus the built-in feed. It runs off Cloudflare Workers and Cloudflare D1 for free.[^1]
3. Brought in comments using [Giscus](https://giscus.app/), which uses the discussions on the blog's [GitHub repo](https://github.com/RdrSeraphim/blog) as a medium. It allows for custom theming, which is one of the more important parts to these integrations for me.
4. Created [blogcli](https://github.com/RdrSeraphim/blogcli) for loose management and format validation of content here. It helps me keep cohesive frontmatter and helps provide for the footnotes system that I took from Ghost – where they are written like `[\^This is a footnote]`. `blogcli` parses them into conventional Markdown footnotes.
5. Pivoted my blog management to CLI. Beyond `blogcli`, the blog is mostly managed in Git. All writings are done in [neovim](https://neovim.io) with a bunch of different plugins, including a custom one for auto-formatting with `blogcli`. I might make a post on my neovim setup sometime, but until then you can check out [my setup repo](https://github.com/RdrSeraphim/kickstart.nvim).

I also expanded [IndieWeb](http://indieweb.org/) integrations, so this blog now supports webmentions alongside rel-me patterns for IndieAuth and RelMeAuth. There's still some work to be done in helping facilitate webmentions, namely with publishing webmentions from this blog and receiving them from post interactions via ssg-ap-bridge, but for now it is mostly stable.

Otherwise, here's the rough stack of this blog (the tl;dr):

* Generator: [Hugo](https://gohugo.io)
* Comments: [Giscus](https://giscus.app)
* ActivityPub Integration: [ssg-ap-bridge, by yours truly](https://github.com/RdrSeraphim/ssg-ap-bridge)
* Host: [Cloudflare Pages](https://pages.dev)
* Domain registrar: [Believe it or not, Cloudflare!](https://cloudflare.com)
* Version control, CI/CD: [GitHub](https://github.com/RdrSeraphim/blog)
* Theme: [srplife_theme, by yours truly](https://github.com/RdrSeraphim/srplife_theme)
* Management: [blogcli](https://github.com/RdrSeraphim/blogcli) and [neovim](https://neovim.io) ([setup](https://github.com/RdrSeraphim/kickstart.nvim))


## The Author

I'm Rdr. Seraphim Pardee, born Elliott, Seraphim Dòmhnullach in Gaelic-speaking contexts. I'm 27 and live in Grand Blanc, Michigan with my wife, Anna. I'm a software engineer for [iA](https://iarx.com/), a pharmacy automation company.

### Contact/Social

* email
    * personal: me@srp.life ([gpg](https://keys.openpgp.org/pks/lookup?op=get&amp;options=mr&amp;search=0x11FD4E775F8802ACD83B6DE5A2BC63FA0276FE65))
    * biblebot: srp@kerygma.digital
* activitypub: [@srp@defcon.social](https://defcon.social/@srp)
    * you can also get posts from this blog by following [@blog@srp.life](https://defcon.social/@blog@srp.life)
* bluesky: [@srp.life](https://bsky.app/profile/srp.life)
    * similarly, you can get posts from this blog by following [@blog.srp.life.ap.brid.gy](https://bsky.app/profile/blog.srp.life.ap.brid.gy)
* twitter/x: [@RdrSeraphim](https://x.com/RdrSeraphim)
* github: [RdrSeraphim](https://github.com/RdrSeraphim)
* gitlab[^2]: [RdrSeraphim](https://gitlab.com/RdrSeraphim)
* linkedin: [srpx](https://linkedin.com/in/srpx)
* lobsters: [srp](https://lobste.rs/u/srp)
* hn: [xsrp](https://news.ycombinator.com/user?id=xsrp)
* discord: srp
* matrix: [@srp:fedora.im](https://matrix.to/#/@srp:fedora.im)
* wikipedia: [Vypr](https://en.wikipedia.org/wiki/User:Vypr)
* urbit: ~fodryn-tidhut
    * i also own ~difsel-sorsul but it goes unused

You might find more on [Keyoxide](https://keyoxide.org/11FD4E775F8802ACD83B6DE5A2BC63FA0276FE65).

---

These are my projects outside of employment, with various levels of importance.

My primary project is [BibleBot](https://biblebot.xyz) ([GitLab](https://gitlab.com/kerygmadigital/biblebot/biblebot)), a Discord bot for referencing the Bible and other Christian resources. It serves over 90,000 servers with a combined audience of over 9 million users. I've founded a nonprofit around the project, [Kerygma Digital](https://kerygma.digital), that seeks to provide this and similar open-source services.

#### Side projects

* rlshim ([github](https://github.com/RdrSeraphim/rlshim)), a lightweight, native shim for launching RuneLite on Linux with Jagex Accounts
* ssg-ap-bridge ([github](https://github.com/RdrSeraphim/ssg-ap-bridge)), a bridge for your statically generated blog and ActivityPub, using your own domain
* canon_law ([github](https://github.com/RdrSeraphim/canon_law)), a web service for referencing Orthodox canon law. It is not yet finished but most that is lacking is the rest of the content.
* goarch_api ([github](https://github.com/RdrSeraphim/goarch_api)), a Python interface for the Greek Orthodox Archdiocese of America's public-facing (but not publicly documented) Chapel API.

#### Old/defunct projects

* typeprint ([github](https://github.com/RdrSeraphim/typeprint)), a program that prints text files char-by-char * made to imitate Fallout 3 terminals and had its 15 minutes of fame by [this reddit post of mine](https://www.reddit.com/r/unixporn/comments/3ab8kn/lilyterm_ive_been_playing_a_lot_of_fallout_3/)
* irc-osric ([github](https://github.com/RdrSeraphim/irc-osric)), an IRC bot developed in Golang to facilitate OSRIC games * fully compliant with OSRIC v1
* osric-cgi ([github](https://github.com/RdrSeraphim/osric-cgi)), a minor fork of kirbyUK's osric-cgi to work with irc-osric

---

#### Other things I d{o,id}

* You can find my dotfiles [here](https://github.com/RdrSeraphim/dottie) and my NixOS configuration [here](https://github.com/RdrSeraphim/nix).
* I do many things on Discord beyond BibleBot.
    * I run Hesychia ([Discord](https://discord.gg/drYwn5v8mw)), a Discord community for Orthodox Christians that serves around ~600, intended to be a level-headed alternative to many turbulent forums. This has granted me opportunities to collaborate with Rdr. Benjamin Cabe of Theoria ([YouTube](https://www.youtube.com/c/Theoria), [Substack](https://theoriatv.substack.com/)) and Rdr. Bojan Teodosijević of Bible Illustrated ([YouTube](https://www.youtube.com/@BibleIllustrated)).
    * I co-founded Christcord, originally "The Christianity Discord Server."
* I occasionally produce music (mainly instrumental hip hop).
* I spend time editing OrthodoxWiki ([profile](https://orthodoxwiki.org/User:Vypr)) and Orthodox-related articles on Wikipedia ([profile](https://en.wikipedia.org/wiki/User:Vypr)).
* I used to consult at AD30, Inc. for [Catena](http://catenabible.com), a Bible application featuring Patristics commentaries and lectionary readings. I incorporated the Greek Orthodox lectionary into the app, daily readings notifications, and daily quotes from saints.

[^1]: This could incur costs if the federation to the account is large enough. I haven't enough time or popularity to get the full breadth of cost analysis.
[^2]: I mainly use GitLab for [BibleBot](https://biblebot.xyz) related things.
