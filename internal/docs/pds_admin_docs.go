package docs

import (
	_ "embed"
	"strings"
)

//go:embed pdsadmin.html
var pdsAdminDocs string

func PDSAdminDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(pdsAdminDocs, "{{.ServerURL}}", serverURL))
}
