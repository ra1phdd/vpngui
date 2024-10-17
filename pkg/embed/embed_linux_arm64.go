//go:build linux && arm64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-linux-arm64
var fs embed.FS

func getFileName() string {
	return "xray-core-linux-arm64"
}
