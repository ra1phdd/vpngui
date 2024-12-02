//go:build windows && 386

package embed

import (
	"embed"
)

//go:embed bin/xray-core-windows-i386.exe
var fsXray embed.FS

//go:embed tun2socks/tun2socks-windows-i386.exe
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-windows-i386.exe"
}

func getFileTun2socks() string {
	return "tun2socks-windows-i386.exe"
}
