package db

import (
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var Conn *sqlx.DB

func Init(DBPath string) error {
	var err error
	Conn, err = sqlx.Open("sqlite3", DBPath)
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
