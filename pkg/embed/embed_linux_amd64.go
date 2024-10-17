//go:build linux && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-linux-amd64
var fs embed.FS

func getFileName() string {
	return "xray-core-linux-amd64"
}
