package docs

import (
	_ "embed"
	"strings"
)

//go:embed orbissociusadmin.html
var orbisSociusAdminDocs string

func OrbisSociusAdminDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(orbisSociusAdminDocs, "{{.ServerURL}}", serverURL))
}
