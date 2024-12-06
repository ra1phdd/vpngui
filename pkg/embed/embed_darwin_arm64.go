//go:build darwin && arm64

package embed

import (
	"embed"
)

//go:embed xray-core/xray-core-darwin-arm64
var fsXray embed.FS

//go:embed tun2socks/tun2socks-darwin-arm64
var fsTun2socks embed.FS

var fsWintun embed.FS

func getFileXray() string {
	return "xray-core-darwin-arm64"
}

func getFileTun2socks() string {
	return "tun2socks-darwin-arm64"
}

func createFileWintun() error {
	return nil
}
