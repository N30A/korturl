package templates

import (
	"embed"
	"html/template"
)

//go:embed *.html partials/*.html
var files embed.FS

var Templates = template.Must(template.ParseFS(files, "*.html", "partials/*.html"))
