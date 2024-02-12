package openapi

import (
	_ "embed"
	"strings"
)

//go:embed pdsadmin-merged.yaml
var pdsAdminSpec string

func PDSAdminSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(pdsAdminSpec, "{{.ServerURL}}", serverURL))
}
