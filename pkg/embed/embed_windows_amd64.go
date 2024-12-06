//go:build windows && amd64

package embed

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
)

//go:embed xray-core/xray-core-windows-amd64.exe
var fsXray embed.FS

//go:embed tun2socks/tun2socks-windows-amd64.exe
var fsTun2socks embed.FS

//go:embed wintun/wintun-amd64.dll
var fsWintun embed.FS

func getFileXray() string {
	return "xray-core-windows-amd64.exe"
}

func getFileTun2socks() string {
	return "tun2socks-windows-amd64.exe"
}

func createFileWintun() error {
	fileContent, err := fsWintun.ReadFile("wintun/wintun-amd64.dll")
	if err != nil {
		fmt.Printf("Ошибка чтения файла из embed.FS: %v\n", err)
		return err
	}

	tempDir := os.TempDir()
	tempFilePath := filepath.Join(tempDir, "wintun.dll")

	if _, err := os.Stat(tempFilePath); err == nil {
		return nil
	}

	err = os.WriteFile(tempFilePath, fileContent, 0644)
	if err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return err
	}

	return nil
}
