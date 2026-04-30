package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sv3r/sv3r.github.io/cmd/site/internal/content"
	"github.com/sv3r/sv3r.github.io/cmd/site/internal/files"
	"github.com/sv3r/sv3r.github.io/cmd/site/internal/render"
)

func main() {
	must(os.RemoveAll("public"))

	must(os.MkdirAll("public", 0o755))
	must(os.MkdirAll(filepath.Join("public", "posts"), 0o755))

	files.Copy("styles.css", "public/styles.css")

	posts := content.LoadPosts("content/posts")

	render.Render("templates/index.gohtml", "public/index.html", render.PageData{
		Title: "sv3r",
		Posts: posts,
	})

	for _, post := range posts {
		outDir := filepath.Join("public/posts", post.Slug)
		must(os.MkdirAll(outDir, 0o755))

		var otherPosts []content.Post
		for _, candidate := range posts {
			if candidate.Slug == post.Slug {
				continue
			}
			otherPosts = append(otherPosts, candidate)
		}

		render.Render("templates/post.gohtml", filepath.Join(outDir, "index.html"), render.PageData{
			Title: post.Title,
			Posts: otherPosts,
			Post:  post,
		})
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
