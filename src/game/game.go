package game

import (
	"github.com/sirupsen/logrus"

	"../db"
	"../models"
	"../util"
)

var log *logrus.Entry = util.LoggerFor("game")

func PlayTurn() {
	log.Println("playing turn")

	_, err := models.AllPlayers()
	if err != nil {
		log.Error("could not get players. skipping turn. %v", err.Error())
		return
	}

	err = incrementXpForIdle()
	if err != nil {
		log.Error("could not incr xp for idle. %v", err.Error())
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
