package repository

import (
	"go.uber.org/zap"
	"vpngui/internal/app/models"
	"vpngui/pkg/db"
	"vpngui/pkg/logger"
)

type SettingsRepository struct{}

func NewSettings() *SettingsRepository {
	return &SettingsRepository{}
}

func (r *SettingsRepository) GetSettings() (models.Settings, error) {
	logger.Debug("Fetching settings")
	var s models.Settings

	err := db.Conn.Get(&s, `SELECT logger_level, autostart, hide_on_startup, language, stats_update_interval FROM settings WHERE id = 1`)
	if err != nil {
		logger.Error("Failed to fetch settings", zap.Error(err))
		return models.Settings{}, err
	}
	logger.Debug("Settings fetched successfully", zap.Any("settings", s))

	return s, nil
}
