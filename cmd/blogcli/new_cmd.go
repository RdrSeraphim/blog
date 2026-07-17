package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/huh"

	"github.com/rdrseraphim/blog/internal/newpost"
)

const defaultContentPostsDir = "content/posts"

func runNew(args []string) error {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	contentDir := fs.String("dir", defaultContentPostsDir, "content posts directory to create the new post under")
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage: blogcli new [-dir content/posts]

Interactively prompts for a new post's metadata and scaffolds a
content/posts/YYYY/MM/DD/slug/index.md page bundle with matching front
matter.
`)
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	var (
		title, tagsRaw, summary, cover, coverAlt string
		draft                                    = true
		slug                                     string
	)

	titleField := huh.NewInput().
		Title("Title").
		Value(&title).
		Validate(func(s string) error {
			if strings.TrimSpace(s) == "" {
				return fmt.Errorf("title is required")
			}
			return nil
		})

	form := huh.NewForm(
		huh.NewGroup(titleField),
		huh.NewGroup(
			huh.NewInput().
				Title("Slug").
				Description("URL path segment, defaults to a slugified title").
				PlaceholderFunc(func() string {
					return newpost.Slugify(title)
				}, &title).
				Value(&slug),
			huh.NewInput().
				Title("Tags").
				Description("comma separated, e.g. Writings, Spiritual Life").
				Value(&tagsRaw),
			huh.NewText().
				Title("Summary").
				Value(&summary),
			huh.NewInput().
				Title("Cover image URL").
				Description("optional").
				Value(&cover),
			huh.NewInput().
				Title("Cover image alt/caption").
				Description("optional, supports markdown links").
				Value(&coverAlt),
			huh.NewConfirm().
				Title("Draft?").
				Value(&draft),
		),
	)

	if err := form.Run(); err != nil {
		return err
	}

	title = strings.TrimSpace(title)
	slug = strings.TrimSpace(slug)
	if slug == "" {
		slug = newpost.Slugify(title)
	} else {
		slug = newpost.Slugify(slug)
	}

	var tags []string
	for _, t := range strings.Split(tagsRaw, ",") {
		t = strings.TrimSpace(t)
		if t != "" {
			tags = append(tags, t)
		}
	}

	p := newpost.Post{
		Title:    title,
		Slug:     slug,
		Tags:     tags,
		Summary:  strings.TrimSpace(summary),
		Cover:    strings.TrimSpace(cover),
		CoverAlt: strings.TrimSpace(coverAlt),
		Draft:    draft,
		Date:     time.Now(),
	}

	path, err := newpost.Create(*contentDir, p)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stderr, "created %s\n", path)
	return nil
}
