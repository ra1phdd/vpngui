package repository

import (
	"database/sql"
	"errors"
	"go.uber.org/zap"
	"vpngui/internal/app/models"
	"vpngui/pkg/db"
	"vpngui/pkg/logger"
)

type RoutesRepository struct{}

func NewRoutes() *RoutesRepository {
	return &RoutesRepository{}
}

func (rr *RoutesRepository) GetRoutes(listType string) (models.ListConfig, error) {
	logger.Debug("Fetching routes", zap.String("listType", listType))
	var listConfig models.ListConfig

	if listType != "blacklist" && listType != "whitelist" {
		logger.Error("Invalid listType", zap.String("listType", listType))
		return models.ListConfig{}, errors.New("listType должен принимать значения 'blacklist' или 'whitelist'")
	}

	err := db.Conn.Get(&listConfig, `SELECT type, domainStrategy, domainMatcher FROM list_config WHERE type = $1`, listType)
	if err != nil {
		logger.Error("Failed to fetch list configuration", zap.String("listType", listType), zap.Error(err))
		return models.ListConfig{}, err
	}

	rows, err := db.Conn.Query(`SELECT rule_type, rule_value FROM rules WHERE list_config_id = (SELECT id FROM list_config WHERE type = $1)`, listType)
	if err != nil {
		logger.Error("Failed to fetch rules for list", zap.String("listType", listType), zap.Error(err))
		return models.ListConfig{}, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			logger.Error("Failed to close rows channel", zap.String("listType", listType), zap.Error(err))
			return
		}
	}(rows)

	for rows.Next() {
		var rule models.Rule
		if err = rows.Scan(&rule.RuleType, &rule.RuleValue); err != nil {
			logger.Error("Failed to scan rule", zap.Error(err))
			return models.ListConfig{}, err
		}
		listConfig.Rules = append(listConfig.Rules, rule)
	}
	logger.Debug("Fetched routes successfully", zap.String("listType", listType), zap.Int("ruleCount", len(listConfig.Rules)))

	return listConfig, nil
}

func (rr *RoutesRepository) AddRule(listType, ruleType, ruleValue string) error {
	logger.Debug("Adding rule", zap.String("listType", listType), zap.String("ruleType", ruleType), zap.String("ruleValue", ruleValue))

	if listType != "blacklist" && listType != "whitelist" {
		logger.Error("Invalid listType", zap.String("listType", listType))
		return errors.New("listType должен принимать значения 'blacklist' или 'whitelist'")
	}

	if ruleType != "domain" && ruleType != "ip" && ruleType != "port" {
		logger.Error("Invalid ruleType", zap.String("ruleType", ruleType))
		return errors.New("ruleType должен принимать значения 'domain', 'ip', или 'port'")
	}

	_, err := db.Conn.Exec(`
		INSERT INTO rules (list_config_id, rule_type, rule_value)
		VALUES ((SELECT id FROM list_config WHERE type = $1), $2, $3)
	`, listType, ruleType, ruleValue)
	if err != nil {
		logger.Error("Failed to add rule to database", zap.Error(err))
		return err
	}
	logger.Debug("Rule added successfully", zap.String("listType", listType), zap.String("ruleType", ruleType), zap.String("ruleValue", ruleValue))

	return nil
}

func (rr *RoutesRepository) DeleteRule(listType, ruleType, ruleValue string) error {
	logger.Debug("Deleting rule", zap.String("listType", listType), zap.String("ruleType", ruleType), zap.String("ruleValue", ruleValue))

	if listType != "blacklist" && listType != "whitelist" {
		logger.Error("Invalid listType", zap.String("listType", listType))
		return errors.New("listType должен принимать значения 'blacklist' или 'whitelist'")
	}

	if ruleType != "domain" && ruleType != "ip" && ruleType != "port" {
		logger.Error("Invalid ruleType", zap.String("ruleType", ruleType))
		return errors.New("ruleType должен принимать значения 'domain', 'ip', или 'port'")
	}

	_, err := db.Conn.Exec(`
		DELETE FROM rules
		WHERE list_config_id = (SELECT id FROM list_config WHERE type = $1)
		AND rule_type = $2 AND rule_value = $3
	`, listType, ruleType, ruleValue)
	if err != nil {
		logger.Error("Failed to delete rule from database", zap.Error(err))
		return err
	}
	logger.Debug("Rule deleted successfully", zap.String("listType", listType), zap.String("ruleType", ruleType), zap.String("ruleValue", ruleValue))

	return nil
}
