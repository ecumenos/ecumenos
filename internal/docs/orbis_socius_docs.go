package docs

import (
	_ "embed"
	"strings"
)

//go:embed orbissocius.html
var orbisSociusDocs string

func OrbisSociusDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(orbisSociusDocs, "{{.ServerURL}}", serverURL))
}
