package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"vpngui/internal/app/models"
)

type Config struct{}

type StructJSON struct {
	ActiveVPN       bool `json:"active-vpn"`
	DisableRoutes   bool `json:"disable-routes"`
	EnableBlackList bool `json:"enable-black-list"`
	EnableWhiteList bool `json:"enable-white-list"`
}

var (
	JSON       StructJSON
	Xray       models.Config
	Routes     models.RoutingConfig
	LockJSON   sync.Mutex
	LockRoutes sync.Mutex
	LockXray   sync.Mutex
	FileJSON   = "config/config.json"
	FileRoutes = "config/routes.json"
	FileXray   = "config/xray.json"
)

func LoadConfig() error {
	if err := Load(FileJSON, &LockJSON, &JSON); err != nil {
		return fmt.Errorf("ошибка загрузки JSON конфига: %s", err)
	}

	if err := Load(FileXray, &LockXray, &Xray); err != nil {
		return fmt.Errorf("ошибка загрузки Xray конфига: %s", err)
	}

	if err := Load(FileRoutes, &LockRoutes, &Routes); err != nil {
		return fmt.Errorf("ошибка загрузки Routes конфига: %s", err)
	}

	return nil
}

func Load(file string, lock *sync.Mutex, config interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	osFile, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	if err = json.Unmarshal(osFile, config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %v", err)
	}

	return nil
}

func SaveConfig() error {
	if err := Save(FileJSON, &LockJSON, &JSON); err != nil {
		return fmt.Errorf("ошибка загрузки JSON конфига: %s", err)
	}

	if err := Save(FileXray, &LockXray, &Xray); err != nil {
		return fmt.Errorf("ошибка загрузки Xray конфига: %s", err)
	}

	if err := Save(FileRoutes, &LockRoutes, &Routes); err != nil {
		return fmt.Errorf("ошибка загрузки Routes конфига: %s", err)
	}

	return nil
}

func Save(file string, lock *sync.Mutex, config interface{}) error {
	lock.Lock()
	defer lock.Unlock()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

func (c *Config) GetJSON() StructJSON {
	return JSON
}

func (c *Config) GetXray() models.Config {
	return Xray
}
