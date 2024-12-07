package main

import (
	"embed"
	"log"
	pkgApp "vpngui/internal/pkg/app"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	err := pkgApp.New(assets)
	if err != nil {
		log.Fatal(err)
	}
}
