//go:build linux && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-linux-amd64
var fsXray embed.FS

//go:embed tun2socks/tun2socks-linux-amd64
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-linux-amd64"
}

func getFileTun2socks() string {
	return "tun2socks-linux-amd64"
}
