//go:build windows && arm64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-windows-arm64.exe
var fs embed.FS

func getFileName() string {
	return "xray-core-windows-arm64.exe"
}
