package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"vpngui/internal/app/models"
)

type Config struct{}

var (
	Xray     models.Xray
	LockXray sync.Mutex
	FileXray = "config/xray.json"
)

func LoadConfig() error {
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

func SaveConfig() error {
	LockXray.Lock()
	defer LockXray.Unlock()

	data, err := json.MarshalIndent(&Xray, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(FileXray, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func (c *Config) GetXray() models.Xray {
	return Xray
}
