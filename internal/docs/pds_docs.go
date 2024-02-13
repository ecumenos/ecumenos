package docs

import (
	_ "embed"
	"strings"
)

//go:embed pds.html
var pdsDocs string

func PDSDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(pdsDocs, "{{.ServerURL}}", serverURL))
}
