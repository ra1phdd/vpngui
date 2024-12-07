//go:build linux && arm64

package embed

import (
	"embed"
)

//go:embed xray-core/xray-core-linux-arm64
var fsXray embed.FS

//go:embed tun2socks/tun2socks-linux-arm64
var fsTun2socks embed.FS

func getFileXray() string {
	return "xray-core-linux-arm64"
}

func getFileTun2socks() string {
	return "tun2socks-linux-arm64"
}

func createFileWintun() error {
	return nil
}
