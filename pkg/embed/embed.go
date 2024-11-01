package embed

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"vpngui/config"
)

//go:embed certs/cert.pem
var embeddedCert []byte

//go:embed certs/key.pem
var embeddedKey []byte

//go:embed config/xray.json
var embeddedXray []byte

//go:embed db/vpngui.db
var embeddedDB []byte

var nameFile = make(map[string]string)

func Init() error {
	Configs()
	Certs()
	DB()

	name := getFileName()

	tempFile, err := os.CreateTemp("", name)
	if err != nil {
		log.Fatalf("Ошибка создания временного файла: %v", err)
		return err
	}

	fileData, err := fs.ReadFile(fmt.Sprintf("bin/%s", name))
	if err != nil {
		log.Fatalf("Ошибка чтения бинарника: %v", err)
		return err
	}
	if _, err := tempFile.Write(fileData); err != nil {
		log.Fatalf("Ошибка записи бинарника Streamlink: %v", err)
		return err
	}
	tempFile.Close()

	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		log.Fatalf("Ошибка установки прав на выполнение: %v", err)
		return err
	}

	nameFile[name] = tempFile.Name()

	return nil
}

func GetTempFileName() string {
	name := getFileName()

	return nameFile[name]
}

func Certs() {
	certDir := "certs"
	FileCert := certDir + "/cert.pem"
	FileKey := certDir + "/key.pem"

	if err := os.MkdirAll(certDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return
	}

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

func Configs() {
	configDir := "config"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return
	}

	if _, err := os.Stat(config.FileXray); os.IsNotExist(err) {
		err = os.WriteFile(config.FileXray, embeddedXray, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}
}

func DB() {
	dbDir := "db"
	FileDB := dbDir + "/vpngui.db"

	if err := os.MkdirAll(dbDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return
	}

	if _, err := os.Stat(FileDB); os.IsNotExist(err) {
		err = os.WriteFile(FileDB, embeddedDB, 0644)
		if err != nil {
			fmt.Printf("Ошибка записи файла: %v\n", err)
			return
		}
	}
}
