//go:build windows && arm64

package embed

import (
	"embed"
)

//go:embed xray-core/xray-core-windows-arm64.exe
var fsXray embed.FS

//go:embed tun2socks/tun2socks-windows-arm64.exe
var fsTun2socks embed.FS

//go:embed wintun/wintun-arm64.dll
var fsWintun embed.FS

func getFileXray() string {
	return "xray-core-windows-arm64.exe"
}

func getFileTun2socks() string {
	return "tun2socks-windows-arm64.exe"
}

func createFileWintun() error {
	fileContent, err := fsWintun.ReadFile("wintun/wintun-arm64.dll")
	if err != nil {
		fmt.Printf("Ошибка чтения файла из embed.FS: %v\n", err)
		return err
	}

	err = os.WriteFile("wintun.dll", fileContent, 0644)
	if err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return err
	}

	return nil
}
