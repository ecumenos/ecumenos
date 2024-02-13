package openapi

import (
	_ "embed"
	"strings"
)

//go:embed orbissocius-merged.yaml
var orbisSociusSpec string

func OrbisSociusSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(orbisSociusSpec, "{{.ServerURL}}", serverURL))
}
