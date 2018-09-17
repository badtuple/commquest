package game

import (
	"fmt"
	"math"
	"math/rand"

	"../db"
	"../frontends"
	"../models"
	"github.com/jmoiron/sqlx"
)

func levelFromXp(xp int) int {
	const base float64 = 0.68
	level := base * math.Cbrt(float64(xp))
	return int(level)
}

func incrementLevels() error {
	ps, err := models.AllPlayers()
	if err != nil {
		return err
	}

	var affected int
	err = db.Transact(func(tx *sqlx.Tx) error {
		for _, p := range ps {
			level := levelFromXp(p.XP)
			if p.Level >= level {
				continue
			}

			// random stat to be incremented
			stat := models.StatList[rand.Intn(len(models.StatList))]

			query := fmt.Sprintf(`
				UPDATE players
				SET level = level + 1, %v = %v + 1, updated_at = NOW()
				WHERE id = $1
				RETURNING xp, level, strength, charisma, intellect, agility, luck
				`, stat, stat,
			)

			err = tx.QueryRowx(query, p.ID).StructScan(&p)
			if err != nil {
				return err
			}

			msg := fmt.Sprintf("%v is now level %v!\n%v",
				p.NameAndTitle(),
				p.Level,
				leveledUpMsg(p),
			)

			frontend.PushMessage(msg)

			affected++
		}
		return nil
	})

	if err != nil {
		return err
	}

	if affected > 0 {
		log.Printf("%v players leveled up", affected)
	}
	return nil
}

func leveledUpMsg(p models.Player) string {
	return fmt.Sprintf(
		"%v has leveled up! *Level*: %v, *XP*: %v, *Strength*: %v, "+
			"*Charisma*: %v, *Intellect*: %v, *Agility*: %v ,*Luck*: %v",

		p.NameAndTitle(),
		p.Level,
		p.XP,
		p.Strength,
		p.Charisma,
		p.Intellect,
		p.Agility,
		p.Luck,
	)
}
