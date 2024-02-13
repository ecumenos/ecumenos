package openapi

import (
	_ "embed"
	"strings"
)

//go:embed pds-merged.yaml
var pdsSpec string

func PDSSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(pdsSpec, "{{.ServerURL}}", serverURL))
}
