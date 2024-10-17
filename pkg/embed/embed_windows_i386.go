//go:build windows && 386

package embed

import (
	"embed"
)

//go:embed bin/xray-core-windows-i386.exe
var fs embed.FS

func getFileName() string {
	return "xray-core-windows-i386.exe"
}
