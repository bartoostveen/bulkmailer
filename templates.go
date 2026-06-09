package templates

import (
	"embed"
	"text/template"
)

// separate package because Go enforces package boundaries

//go:embed *.txt
var templateFS embed.FS
var Templates = template.Must(template.ParseFS(templateFS, "*.txt"))
