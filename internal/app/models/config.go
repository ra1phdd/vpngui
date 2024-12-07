package models

type Config struct {
	ActiveVPN     bool   `db:"active_vpn"`
	DisableRoutes bool   `db:"disable_routes"`
	ListMode      string `db:"list_mode"` // Принимает значения "blacklist" или "whitelist"
	VpnMode       string `db:"vpn_mode"`  // Принимает значения "tun" или "proxy"
}
