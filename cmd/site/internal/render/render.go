package render

import (
	"html/template"
	"log"
	"os"

	"github.com/sv3r/sv3r.github.io/cmd/site/internal/content"
)

type PageData struct {
	Title string
	Posts []content.Post
	Post  content.Post
}

func Render(templatePath, outputPath string, data PageData) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.gohtml",
		templatePath,
	))

	file, err := os.Create(outputPath)
	must(err)

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	must(tmpl.ExecuteTemplate(file, "base", data))
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
