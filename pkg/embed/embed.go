package embed

import (
	"embed"
	_ "embed"
	"fmt"
	"log"
	"os"
	"vpngui/internal/app/config"
)

//go:embed config/xray.json
var embeddedXray []byte

//go:embed db/vpngui.db
var embeddedDB []byte

var nameFile = make(map[string]string)

func Init() error {
	Configs()
	DB()

	err := CreateFile(getFileXray(), "bin", fsXray)
	if err != nil {
		return err
	}
	err = CreateFile(getFileTun2socks(), "tun2socks", fsTun2socks)
	if err != nil {
		return err
	}

	return nil
}

func CreateFile(name string, dir string, fs embed.FS) error {
	tempFile, err := os.CreateTemp("", name)
	if err != nil {
		log.Fatalf("Ошибка создания временного файла: %v", err)
		return err
	}

	fileData, err := fs.ReadFile(fmt.Sprintf("%s/%s", dir, name))
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

func GetTempFileName(name string) string {
	switch name {
	case "xray-core":
		return nameFile[getFileXray()]
	case "tun2socks":
		return nameFile[getFileTun2socks()]
	}

	return ""
}

func Configs() {
	configDir := "config"

	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Printf("Ошибка создания директории: %v\n", err)
		return
	}

	err := os.WriteFile(config.FileXray, embeddedXray, 0644)
	if err != nil {
		fmt.Printf("Ошибка записи файла: %v\n", err)
		return
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
