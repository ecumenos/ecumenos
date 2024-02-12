package openapi

import (
	_ "embed"
	"strings"
)

//go:embed orbissociusadmin-merged.yaml
var orbisSociusAdminSpec string

func OrbisSociusAdminSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(orbisSociusAdminSpec, "{{.ServerURL}}", serverURL))
}
