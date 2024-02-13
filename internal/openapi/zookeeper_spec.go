package openapi

import (
	_ "embed"
	"strings"
)

//go:embed zookeeper-merged.yaml
var zookeeperSpec string

func ZookeeperSpec(serverURL string) []byte {
	return []byte(strings.ReplaceAll(zookeeperSpec, "{{.ServerURL}}", serverURL))
}
