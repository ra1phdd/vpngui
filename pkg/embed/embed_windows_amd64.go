//go:build windows && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-windows-amd64.exe
var fs embed.FS

func getFileName() string {
	return "xray-core-windows-amd64.exe"
}
