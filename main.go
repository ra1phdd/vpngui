package main

import (
	"embed"
	"log"
	pkgApp "vpngui/internal/pkg/app"
	embeds "vpngui/pkg/embed"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	embeds.InitConfigs()
	embeds.InitCerts()

	err := pkgApp.New(assets)
	if err != nil {
		log.Fatal(err)
	}
}
