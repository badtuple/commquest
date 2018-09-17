package models

import "../db"

type Item struct {
	ID int `db:"id"`

	Name    string `db:"name"`
	Article string `db:"article"`

	XPIncr        int `db:"xp_incr"`
	LevelIncr     int `db:"level_incr"`
	StrengthIncr  int `db:"strength_incr"`
	CharismaIncr  int `db:"charisma_incr"`
	IntellectIncr int `db:"intellect_incr"`
	AgilityIncr   int `db:"agility_incr"`
	LuckIncr      int `db:"luck_incr"`
}

func CreateItem(name, article string, xp, lvl, s, c, i, a, l int) (Item, error) {
	var item Item
	err := db.PSQL().QueryRowx(`
		INSERT INTO items (
			name, article,
			xp_incr, level_incr,
			strength_incr, charisma_incr,
			intellect_incr, agility_incr,
			luck_incr
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, article,
			xp_incr, level_incr,
			strength_incr, charisma_incr,
			intellect_incr, agility_incr,
			luck_incr
	`, name, article, xp, lvl, s, c, i, a, l).
		StructScan(&item)

	return item, err
}

func AllItems() ([]Item, error) {
	var is []Item
	err := db.PSQL().Select(&is, `
		SELECT id, name, article,
			xp_incr, level_incr,
			strength_incr, charisma_incr,
			intellect_incr, agility_incr,
			luck_incr
		FROM items

	`)
	return is, err
}

func GetRandomItem() (*Item, error) {
	var item Item
	err := db.PSQL().QueryRowx(`
		SELECT id, name, article
		FROM items
		ORDER BY RANDOM()
		LIMIT 1
	`).StructScan(&item)
	return &item, err
}
