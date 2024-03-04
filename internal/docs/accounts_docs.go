package docs

import (
	_ "embed"
	"strings"
)

//go:embed accounts.html
var accountsDocs string

func AccountsDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(accountsDocs, "{{.ServerURL}}", serverURL))
}
