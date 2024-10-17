//go:build darwin && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-darwin-amd64
var fs embed.FS

func getFileName() string {
	return "xray-core-darwin-amd64"
}
