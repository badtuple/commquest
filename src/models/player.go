package models

import (
	"time"

	"../db"
)

type Player struct {
	ID     int    `db:"id" json:"id"`
	Handle string `db:"handle" json:"handle"` // Handle in Frontend

	Name  string `db:"name" json:"name"`   // Name of character
	Class string `db:"class" json:"class"` // Class of character
	XP    int    `db:"xp" json:"xp"`       // Experience points

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

func CreatePlayer(handle, name, class string) (Player, error) {
	var p Player
	err := db.PSQL().QueryRowx(`
		INSERT INTO players (handle, name, class, xp, created_at, updated_at)
		VALUES ($1, $2, $3, 0, NOW(), NOW())
		RETURNING id, handle, name, class, xp, created_at, updated_at
	`, handle, name, class).StructScan(&p)

	return p, err
}

func FindPlayer(id int) (Player, error) {
	var p Player
	err := db.PSQL().QueryRowx(`
		SELECT id, handle, name, class, xp, created_at, updated_at
		FROM players
		WHERE id = $1
		LIMIT 1
	`, id).StructScan(&p)

	return p, err
}

func FindPlayerByHandle(handle string) (Player, error) {
	var p Player
	err := db.PSQL().QueryRowx(`
		SELECT id, handle, name, class, xp, created_at, updated_at
		FROM players
		WHERE handle = $1
		LIMIT 1
	`, handle).StructScan(&p)

	return p, err
}
