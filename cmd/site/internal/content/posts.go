package content

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
)

type Post struct {
	Title       string    `yaml:"title"`
	Date        time.Time `yaml:"date"`
	Author      string    `yaml:"author"`
	Description string    `yaml:"description"`
	Slug        string    `yaml:"slug"`
	Content     template.HTML
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		highlighting.NewHighlighting(
			highlighting.WithStyle("github-dark"),
		),
	),
)

func LoadPosts(dir string) []Post {
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

	if !strings.HasPrefix(text, "---\n") && !strings.HasPrefix(text, "---\r\n") {
		log.Fatal("post is missing frontmatter")
	}

	var post Post
	body, err := frontmatter.Parse(bytes.NewReader(raw), &post)
	if err != nil {
		log.Fatal("invalid frontmatter")
	}

	var buf bytes.Buffer
	must(md.Convert(bytes.TrimSpace(body), &buf))
	post.Content = template.HTML(buf.String())

	return post
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
