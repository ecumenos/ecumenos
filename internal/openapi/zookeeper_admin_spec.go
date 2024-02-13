package openapi

import (
	_ "embed"
	"strings"
)

//go:embed zookeeperadmin-merged.yaml
var zookeeperAdminSpec string

func ZookeeperAdminSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(zookeeperAdminSpec, "{{.ServerURL}}", serverURL))
}
