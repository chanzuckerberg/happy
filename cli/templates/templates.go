package templates

import (
	"embed"
)

//go:embed templates/*.tmpl
var static embed.FS

func StaticAsset(name string) ([]byte, error) {
	return static.ReadFile("templates/" + name)
}
