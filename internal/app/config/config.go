package config

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sync"
	"vpngui/internal/app/models"
	"vpngui/pkg/logger"
)

var (
	Xray     models.Xray
	LockXray sync.Mutex
	FileXray = "config/xray.json"
)

type Config struct{}

func New() *Config {
	err := Load()
	if err != nil {
		logger.Error("Error loading config", zap.Error(err))
	}

	return &Config{}
}

func Load() error {
	LockXray.Lock()
	defer LockXray.Unlock()

	osFile, err := os.ReadFile(FileXray)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err = json.Unmarshal(osFile, &Xray); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}

func Save() error {
	LockXray.Lock()
	defer LockXray.Unlock()

	data, err := json.MarshalIndent(&Xray, "", "  ")
	if err != nil {
		logger.Error("Ошибка преобразования структуры в JSON-конфиг", zap.String("file", FileXray))
		return err
	}

	err = os.WriteFile(FileXray, data, 0644)
	if err != nil {
		logger.Error("Ошибка записи JSON-конфига", zap.String("file", FileXray))
		return err
	}

	return nil
}

func (c *Config) Get() models.Xray {
	return Xray
}
