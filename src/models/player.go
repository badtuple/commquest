package models

import (
	"time"

	"../db"
)

///////////////////////////
// The main Player stuct //
///////////////////////////
//
// This handles both User information and Character information.  These may be
// broken out into different pieces at some point.
//
type Player struct {
	ID     int    `db:"id" json:"id"`
	Handle string `db:"handle" json:"handle"` // Handle in Frontend

	Name  string `db:"name" json:"name"`   // Name of character
	Class string `db:"class" json:"class"` // Class of character
	XP    int    `db:"xp" json:"xp"`       // Experience points
	Level int    `db:"level" json:"level"`

	///////////
	// Stats //
	///////////
	//
	// While these may be used to determine the outcome of fights in the
	// future, the primary purpose of the stats is to determine character
	// choices within quests.
	//
	Strength  int `db:"strength" json:"strength"`
	Charisma  int `db:"charisma" json:"charisma"`
	Intellect int `db:"intellect" json:"intellect"`
	Agility   int `db:"agility" json:"agility"`
	Luck      int `db:"luck" json:"luck"`

	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

var StatList = []string{"strength", "charisma", "intellect", "agility", "luck"}

func (p Player) StatMap() map[string]int {
	return map[string]int{
		"strength":  p.Strength,
		"charisma":  p.Charisma,
		"intellect": p.Intellect,
		"agility":   p.Agility,
		"luck":      p.Luck,
	}
}

// Returns the highest stat for a particular player and it's value
func (p Player) PrimaryStat() (string, int) {
	var stat string
	var val int

	for k, v := range p.StatMap() {
		if v > val {
			stat = k
			val = v
		}
	}

	return stat, val
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
		SELECT
			id, handle, name, class, xp, level,
			strength, charisma, intellect, agility, luck,
			created_at, updated_at
		FROM players
		WHERE id = $1
		LIMIT 1
	`, id).StructScan(&p)

	return p, err
}

func FindPlayerByHandle(handle string) (Player, error) {
	var p Player
	err := db.PSQL().QueryRowx(`
		SELECT
			id, handle, name, class, xp, level,
			strength, charisma, intellect, agility, luck,
			created_at, updated_at
		FROM players
		WHERE handle = $1
		LIMIT 1
	`, handle).StructScan(&p)

	return p, err
}

func AllPlayers() ([]Player, error) {
	var ps []Player
	err := db.PSQL().Select(&ps, `
		SELECT
			id, handle, name, class, xp, level,
			strength, charisma, intellect, agility, luck,
			created_at, updated_at
		FROM players
	`)

	return ps, err
}
