package models

type ListConfig struct {
	Type           string `db:"type"` // Принимает значения "blacklist" или "whitelist"
	Rules          []Rule
	DomainStrategy string `db:"domainStrategy"`
	DomainMatcher  string `db:"domainMatcher"`
}

type Rule struct {
	RuleType  string `db:"rule_type"` // Принимает значения "domain", "ip", или "port"
	RuleValue string `db:"rule_value"`
}
