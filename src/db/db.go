package db

import (
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var onceDB sync.Once
var conn *sqlx.DB

func PSQL() *sqlx.DB {
	onceDB.Do(setup)
	return conn
}

func setup() {
	dbstring := `user=postgres dbname=commquest_development sslmode=disable`
	db, err := sqlx.Open("postgres", dbstring)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	conn = db
}

func Transact(fn func(*sqlx.Tx) error) error {
	tx, err := PSQL().Beginx()
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
