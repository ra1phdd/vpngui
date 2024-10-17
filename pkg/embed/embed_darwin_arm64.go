//go:build darwin && arm64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-darwin-arm64
var fs embed.FS

func getFileName() string {
	return "xray-core-darwin-arm64"
}
