//go:build linux && 386

package embed

import (
	"embed"
)

//go:embed bin/xray-core-linux-i386
var fsXray embed.FS

//go:embed tun2socks/tun2socks-linux-i386
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-linux-i386"
}

func getFileTun2socks() string {
	return "tun2socks-linux-i386"
}
