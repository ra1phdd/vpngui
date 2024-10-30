package models

type Config struct {
	ActiveVPN     bool   `db:"active_vpn"`
	DisableRoutes bool   `db:"disable_routes"`
	ListMode      string `db:"list_mode"` // Принимает значения "blacklist" или "whitelist"
}
