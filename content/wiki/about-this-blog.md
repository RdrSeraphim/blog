+++
date = 2026-08-04T13:18:00-04:00
lastmod = 2026-08-04T13:18:00-04:00
title = 'About This Blog'
draft = false
slug = 'about-this-blog'
category = 'Meta'
summary = 'This blog is a public-facing outlet for my various thoughts and tangents, including some technical postings going over software development or homelabbing.'
+++

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


