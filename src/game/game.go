package game

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	"../db"
	"../models"
	"../util"
)

var log *logrus.Entry = util.LoggerFor("game")

func PlayTurn() {
	log.Println("playing turn")

	err := incrementXpForIdle()
	if err != nil {
		log.Error("could not incr xp for idle. %v", err.Error())
	}

	err = incrementLevels()
	if err != nil {
		log.Error("could not increment level. %v", err.Error())
	}

	log.Println("ending turn")
}

// Right now this increments 1 xp for every turn for every player.
// After an "active" status for the player is implemented, there will be a
// config option to only increment xp for people who are active/in channel.
func incrementXpForIdle() error {
	// TODO: impl status for player and config to only incr active

	result, err := db.PSQL().Exec(`UPDATE players SET xp = xp + 1`)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	log.Printf("incremnted xp for %v players", affected)
	return nil
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

			// Level up
			p.Level = level
			_, err = tx.Exec(`
				UPDATE players
				SET level = level + 1, updated_at = NOW()
				WHERE id = $1
			`, p.ID)

			if err != nil {
				return err
			}

			affected++
		}
		return nil
	})

	if err != nil {
		return err
	}
	log.Printf("incremented level for %v players", affected)
	return nil
}
