package models

type Settings struct {
	LoggerLevel         string `db:"logger_level"`
	Autostart           bool   `db:"autostart"`
	HideOnStartup       bool   `db:"hide_on_startup"`
	Language            string `db:"language"`
	StatsUpdateInterval int    `db:"stats_update_interval"`
}
