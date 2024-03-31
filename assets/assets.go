package assets

import "embed"

//go:embed templates/*.tmpl
var TemplateFS embed.FS
