package embed

import (
	_ "embed"
	"fmt"
	"os"
	"vpngui/config"
)

//go:embed config/config.json
var embeddedJSON []byte

//go:embed config/xray.json
var embeddedXray []byte

//go:embed config/routes.json
var embeddedRoutes []byte

func InitConfigs() {
	if _, err := os.Stat(config.FileJSON); os.IsNotExist(err) {
		err = os.WriteFile(config.FileJSON, embeddedJSON, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}

	if _, err := os.Stat(config.FileXray); os.IsNotExist(err) {
		err = os.WriteFile(config.FileXray, embeddedXray, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}

	if _, err := os.Stat(config.FileRoutes); os.IsNotExist(err) {
		err = os.WriteFile(config.FileRoutes, embeddedRoutes, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}
}
