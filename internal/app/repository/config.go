package repository

import (
	"errors"
	"go.uber.org/zap"
	"vpngui/internal/app/models"
	"vpngui/pkg/db"
	"vpngui/pkg/logger"
)

type ConfigRepository struct{}

func NewConfig() *ConfigRepository {
	return &ConfigRepository{}
}

func (cr *ConfigRepository) GetConfig() (models.Config, error) {
	logger.Debug("Fetching config")
	var c models.Config

	err := db.Conn.Get(&c, `SELECT active_vpn, disable_routes, list_mode, vpn_mode FROM config WHERE id = 1`)
	if err != nil {
		logger.Error("Failed to fetch config", zap.Error(err))
		return models.Config{}, err
	}
	logger.Debug("Config fetched successfully", zap.Any("config", c))

	return c, nil
}

func (cr *ConfigRepository) UpdateActiveVPN(activeVPN bool) error {
	logger.Debug("Updating active VPN status", zap.Bool("activeVPN", activeVPN))

	_, err := db.Conn.Exec(`UPDATE config SET active_vpn = $1 WHERE id = 1`, activeVPN)
	if err != nil {
		logger.Error("Failed to update active VPN status", zap.Error(err))
		return err
	}
	logger.Debug("Active VPN status updated successfully", zap.Bool("activeVPN", activeVPN))

	return nil
}

func (cr *ConfigRepository) UpdateDisableRoutes(disableRoutes bool) error {
	logger.Debug("Updating disable routes status", zap.Bool("disableRoutes", disableRoutes))

	_, err := db.Conn.Exec(`UPDATE config SET disable_routes = $1 WHERE id = 1`, disableRoutes)
	if err != nil {
		logger.Error("Failed to update disable routes status", zap.Error(err))
		return err
	}
	logger.Debug("Disable routes status updated successfully", zap.Bool("disableRoutes", disableRoutes))

	return nil
}

func (cr *ConfigRepository) UpdateListMode(listMode string) error {
	logger.Debug("Updating list mode", zap.String("listMode", listMode))

	if listMode != "blacklist" && listMode != "whitelist" {
		logger.Error("Invalid list mode", zap.String("listMode", listMode))
		return errors.New("list_mode должен принимать значения 'blacklist' или 'whitelist'")
	}

	_, err := db.Conn.Exec(`UPDATE config SET list_mode = $1 WHERE id = 1`, listMode)
	if err != nil {
		logger.Error("Failed to update list mode", zap.Error(err))
		return err
	}
	logger.Debug("List mode updated successfully", zap.String("listMode", listMode))

	return nil
}

func (cr *ConfigRepository) UpdateVPNMode(vpnMode string) error {
	logger.Debug("Updating VPN mode", zap.String("vpnMode", vpnMode))

	if vpnMode != "tun" && vpnMode != "proxy" {
		logger.Error("Invalid VPN mode", zap.String("vpnMode", vpnMode))
		return errors.New("vpn_mode должен принимать значения 'tun' или 'proxy'")
	}

	_, err := db.Conn.Exec(`UPDATE config SET vpn_mode = $1 WHERE id = 1`, vpnMode)
	if err != nil {
		logger.Error("Failed to update VPN mode", zap.Error(err))
		return err
	}
	logger.Debug("VPN mode updated successfully", zap.String("vpnMode", vpnMode))

	return nil
}
