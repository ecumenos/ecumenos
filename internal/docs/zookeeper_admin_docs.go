package docs

import (
	_ "embed"
	"strings"
)

//go:embed zookeeperadmin.html
var zookeeperAdminDocs string

func ZookeeperAdminDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(zookeeperAdminDocs, "{{.ServerURL}}", serverURL))
}
