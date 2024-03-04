package openapi

import (
	_ "embed"
	"strings"
)

//go:embed accounts-merged.yaml
var accountsSpec string

func AccountsSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(accountsSpec, "{{.ServerURL}}", serverURL))
}
