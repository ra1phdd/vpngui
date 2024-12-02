//go:build darwin && amd64

package embed

import (
	"embed"
)

//go:embed bin/xray-core-darwin-amd64
var fsXray embed.FS

//go:embed tun2socks/tun2socks-darwin-amd64
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-darwin-amd64"
}

func getFileTun2socks() string {
	return "tun2socks-darwin-amd64"
}
