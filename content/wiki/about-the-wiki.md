+++
date = 2026-08-30T21:46:50-04:00
lastmod = 2026-08-30T21:46:50-04:00
title = 'About The Wiki'
draft = false
slug = 'about-the-wiki'
category = 'Meta'
summary = "This is my take on a flat-file wiki by leveraging Hugo's capabilities. With some extras, I can dedicate a whole section of my blog to elaborating on certain topics, referencing them amongst wiki entries and in the rest of the blog."
+++

This is my take on a flat-file wiki by leveraging Hugo's capabilities. With some extras, I can dedicate a whole section of my blog to elaborating on certain topics, referencing them amongst wiki entries and in the rest of the blog.

### Inspirations

Like many elements to my blog, this is inspired in part by the backlinking (and correlating "popups") feature on [Gwern Branwen's blog](https://gwern.net/design). It's also inspired by Alex Schroeder's [Oddμ](https://src.alexschroeder.ch/oddmu.git/), which he uses to facilitate his [blog](https://alexschroeder.ch). I really liked the idea of having a statically generated wiki in a manner that was simple enough to not carry so much overhead in running the thing and to just focus on writing, thus it felt easier to make a wiki embedded with my blog vs. self-hosting oddμ on a separate domain.

### Usage and Implementation

Whereas [my journal](https://srp.life/on-the-keeping-of-a-journal/) exists for private reflection and thoughts, this blog is a public-facing analogue. This wiki, in turn, allows for me to write up entries about isolated topics I might be nerdy about, which also allows me to build context to terms and concepts that might not be familiar to my readers. Thanks to the wiki referencing feature I added, I can easily link to a wiki entry where a summary can be presented on a concept simply by hovering over the text.

So, my intention in usage is largely to establish my own form of reference material for others so that my blog entries might be a bit more vivid.

As an example, I can easily reference a wiki page like the {{< wikiref "oca" >}}OCA{{< /wikiref >}}. This example in the source looks like:
```markdown
{{</* wikiref "oca" */>}}Orthodox Church in America{{</* /wikiref */>}}
```

I can also do like `{{</* wikiref "oca" /*/>}}` and get {{< wikiref "oca" />}}.

These wikirefs are half [Hugo shortcode](https://github.com/RdrSeraphim/srplife_theme/blob/modern/layouts/_shortcodes/wikiref.html) and half [JavaScript](https://github.com/RdrSeraphim/srplife_theme/blob/modern/static/js/wikiref.js) for making sure the tooltips align properly.
