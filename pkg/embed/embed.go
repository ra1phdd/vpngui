package embed

import (
	"fmt"
	"log"
	"os"
)

var nameFile = make(map[string]string)

func Init() error {
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
