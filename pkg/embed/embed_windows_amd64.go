//go:build windows && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-windows-amd64.exe
var fsXray embed.FS

//go:embed tun2socks/tun2socks-windows-amd64.exe
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-windows-amd64.exe"
}

func getFileTun2socks() string {
	return "tun2socks-windows-amd64.exe"
}
