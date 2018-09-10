package main

import (
	"log"

	"../src/models"
)

func main() {
	log.Println("seeding database")

	log.Println("creating player")
	p, err := models.CreatePlayer("mdcox", "mdcox", "developer")
	if err != nil {
		panic(err)
	}

	_p, err := models.FindPlayer(p.ID)
	if err != nil {
		panic(err)
	}

	if p.ID != _p.ID {
		panic("inserted p.ID does not match retrieved _p.ID")
	}

	log.Printf("created player: %+v", p)

	log.Println("finished seeding")
}
