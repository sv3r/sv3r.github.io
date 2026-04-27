package main

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

type Post struct {
	Title       string
	Date        time.Time
	Author      string
	Description string
	Slug        string
	Content     template.HTML
}

type PageData struct {
	Title string
	Posts []Post
	Post  Post
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("github-dark"),
		),
	),
)

func main() {
	must(os.RemoveAll("public"))

	must(os.MkdirAll("public", 0o755))
	must(os.MkdirAll(filepath.Join("public", "posts"), 0o755))

	copyFile("styles.css", "public/styles.css")

	posts := loadPosts("content/posts")

	render("templates/index.gohtml", "public/index.html", PageData{
		Title: "sv3r",
		Posts: posts,
	})

	for _, post := range posts {
		outDir := filepath.Join("public/posts", post.Slug)
		must(os.MkdirAll(outDir, 0o755))

		var otherPosts []Post
		for _, candidate := range posts {
			if candidate.Slug == post.Slug {
				continue
			}
			otherPosts = append(otherPosts, candidate)
		}

		render("templates/post.gohtml", filepath.Join(outDir, "index.html"), PageData{
			Title: post.Title,
			Posts: otherPosts,
			Post:  post,
		})
	}
}

func loadPosts(dir string) []Post {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		log.Fatal(err)
	}

	var posts []Post

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		raw, err := os.ReadFile(path)
		must(err)

		post := parsePost(raw)
		if post.Slug == "" {
			post.Slug = strings.TrimSuffix(entry.Name(), ".md")
		}

		posts = append(posts, post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Date.After(posts[j].Date)
	})

	return posts
}

func parsePost(raw []byte) Post {
	text := string(raw)

	if !strings.HasPrefix(text, "---\n") {
		log.Fatal("post is missing frontmatter")
	}

	parts := strings.SplitN(text, "---\n", 3)
	if len(parts) < 3 {
		log.Fatal("invalid frontmatter")
	}

	meta := strings.TrimSpace(parts[1])
	body := strings.TrimSpace(parts[2])

	post := Post{}

	for _, line := range strings.Split(meta, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}

		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"`)

		switch key {
		case "title":
			post.Title = value
		case "author":
			post.Author = value
		case "description":
			post.Description = value
		case "slug":
			post.Slug = value
		case "date":
			date, err := time.Parse("2006-01-02", value)
			must(err)
			post.Date = date
		}
	}

	var buf bytes.Buffer
	must(md.Convert([]byte(body), &buf))
	post.Content = template.HTML(buf.String())

	return post
}

func render(templatePath, outputPath string, data PageData) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.gohtml",
		templatePath,
	))

	file, err := os.Create(outputPath)
	must(err)
	defer file.Close()

	must(tmpl.ExecuteTemplate(file, "base", data))
}

func copyFile(src, dst string) {
	input, err := os.ReadFile(src)
	must(err)

	must(os.WriteFile(dst, input, 0o644))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
