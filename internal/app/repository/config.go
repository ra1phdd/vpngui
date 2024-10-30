package repository

import (
	"errors"
	"vpngui/internal/app/models"
	"vpngui/pkg/db"
)

type ConfigRepository struct{}

func NewConfig() *ConfigRepository {
	return &ConfigRepository{}
}

func (r *ConfigRepository) GetConfig() (models.Config, error) {
	var c models.Config

	err := db.Conn.Get(&c, `SELECT active_vpn, disable_routes, list_mode FROM config WHERE id = 1`)
	return c, err
}

func (r *ConfigRepository) UpdateActiveVPN(activeVPN bool) error {
	_, err := db.Conn.Exec(`UPDATE config SET active_vpn = $1 WHERE id = 1`, activeVPN)
	return err
}

func (r *ConfigRepository) UpdateDisableRoutes(disableRoutes bool) error {
	_, err := db.Conn.Exec(`UPDATE config SET disable_routes = $1 WHERE id = 1`, disableRoutes)
	return err
}

func (r *ConfigRepository) UpdateListMode(listMode string) error {
	if listMode != "blacklist" && listMode != "whitelist" {
		return errors.New("list_mode must be either 'blacklist' or 'whitelist'")
	}

	_, err := db.Conn.Exec(`UPDATE config SET list_mode = ? WHERE id = 1`, listMode)
	return err
}
