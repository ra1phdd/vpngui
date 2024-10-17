package embed

import (
	_ "embed"
	"fmt"
	"os"
)

//go:embed certs/cert.pem
var embeddedCert []byte

//go:embed certs/key.pem
var embeddedKey []byte

func InitCerts() {
	FileCert := "config/cert.pem"
	FileKey := "config/key.pem"

	if _, err := os.Stat(FileCert); os.IsNotExist(err) {
		err = os.WriteFile(FileCert, embeddedCert, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}

	if _, err := os.Stat(FileKey); os.IsNotExist(err) {
		err = os.WriteFile(FileKey, embeddedKey, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}
}
