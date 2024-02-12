package docs

import (
	_ "embed"
	"strings"
)

//go:embed zookeeper.html
var zookeeperDocs string

func ZookeeperDocs(serverURL string) []byte {
	return []byte(strings.ReplaceAll(zookeeperDocs, "{{.ServerURL}}", serverURL))
}
