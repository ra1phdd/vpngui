package repository

import (
	"errors"
	"vpngui/internal/app/models"
	"vpngui/pkg/db"
)

type RoutesRepository struct{}

func NewRoutes() *RoutesRepository {
	return &RoutesRepository{}
}

func (r *RoutesRepository) GetRoutes(listType string) (models.ListConfig, error) {
	var listConfig models.ListConfig

	if listType != "blacklist" && listType != "whitelist" {
		return listConfig, errors.New("listType must be either 'blacklist' or 'whitelist'")
	}

	err := db.Conn.Get(&listConfig, `SELECT type, domainStrategy, domainMatcher FROM list_config WHERE type = $1`, listType)
	if err != nil {
		return listConfig, err
	}

	rows, err := db.Conn.Query(`SELECT rule_type, rule_value FROM rules WHERE list_config_id = (SELECT id FROM list_config WHERE type = $1)`, listType)
	if err != nil {
		return listConfig, err
	}
	defer rows.Close()

	for rows.Next() {
		var rule models.Rule
		if err := rows.Scan(&rule.RuleType, &rule.RuleValue); err != nil {
			return listConfig, err
		}
		listConfig.Rules = append(listConfig.Rules, rule)
	}

	return listConfig, nil
}

func (r *RoutesRepository) AddRule(listType, ruleType, ruleValue string) error {
	if listType != "blacklist" && listType != "whitelist" {
		return errors.New("listType must be either 'blacklist' or 'whitelist'")
	}

	if ruleType != "domain" && ruleType != "ip" && ruleType != "port" {
		return errors.New("ruleType must be either 'domain', 'ip', or 'port'")
	}

	_, err := db.Conn.Exec(`
		INSERT INTO rules (list_config_id, rule_type, rule_value)
		VALUES ((SELECT id FROM list_config WHERE type = $1), $2, $3)
	`, listType, ruleType, ruleValue)

	return err
}

func (r *RoutesRepository) DeleteRule(listType, ruleType, ruleValue string) error {
	if listType != "blacklist" && listType != "whitelist" {
		return errors.New("listType must be either 'blacklist' or 'whitelist'")
	}

	if ruleType != "domain" && ruleType != "ip" && ruleType != "port" {
		return errors.New("ruleType must be either 'domain', 'ip', or 'port'")
	}

	_, err := db.Conn.Exec(`
		DELETE FROM rules
		WHERE list_config_id = (SELECT id FROM list_config WHERE type = $1)
		AND rule_type = $2 AND rule_value = $3
	`, listType, ruleType, ruleValue)

	return err
}
