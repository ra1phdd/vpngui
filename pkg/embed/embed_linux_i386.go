//go:build linux && 386

package embed

import (
	"embed"
)

//go:embed bin/xray-core-linux-i386
var fs embed.FS

func getFileName() string {
	return "xray-core-linux-i386"
}
