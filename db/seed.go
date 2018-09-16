package main

import (
	"log"

	"../src/models"
)

func main() {
	log.Println("seeding database")

	log.Println("creating player")
	createPlayer()

	log.Print("loading seed items")
	loadItems()

	log.Println("finished seeding")
}

func createPlayer() {
	p, err := models.CreatePlayer("badtuple", "badtuple", "developer")
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
}

func loadItems() {
	items := []models.Item{
		{0, "Eye of Saggathah", "the", 0, 0, 0, 0, 5, 0, 1},
		{0, "punic artifact", "a", 0, 0, 2, 0, 0, 1, 0},
		{0, "Cleopatras Crown", "", 0, 0, 0, 0, 6, 0, 1},
		{0, "Western Wind", "the", 0, 0, 0, 0, 0, 3, 0},
		{0, "lucky dice", "some", 0, 0, 0, 0, 0, 0, 3},
		{0, "dice, with smiley faces carved into them", "some", 0, 0, 0, 0, 0, 0, 4},
		{0, "thieves gloves", "some", 0, 0, 0, 0, 0, 3, 0},
		{0, "pet rock", "a", 100, 0, 0, 0, -1, 0, 0},
		{0, "Goliath's Belt", "", 0, 0, 5, 0, 0, 0, 0},
	}

	for _, i := range items {
		_, err := models.CreateItem(i.Name, i.Article, i.XPIncr, i.LevelIncr, i.StrengthIncr, i.CharismaIncr, i.IntellectIncr, i.AgilityIncr, i.LuckIncr)
		if err != nil {
			panic(err.Error())
		}
	}

}
