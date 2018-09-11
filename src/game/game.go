package game

import (
	"../util"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("game")

func PlayTurn() {
	log.Println("playing turn")

	log.Println("ending turn")
}
