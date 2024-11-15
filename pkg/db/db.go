package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

var Conn *sqlx.DB

func Init(DBPath string) error {
	var err error
	Conn, err = sqlx.Open("sqlite", DBPath)
	if err != nil {
		return err
	}

	Conn.SetMaxOpenConns(1)
	Conn.SetMaxIdleConns(1)
	Conn.SetConnMaxLifetime(time.Hour)

	err = Conn.Ping()
	if err != nil {
		return err
	}

	return nil
}
