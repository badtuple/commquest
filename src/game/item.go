package game

import (
	"fmt"
	"math/rand"
	"strings"

	"../frontends"
	"../models"
)

func maybeDropItem() error {
	// Drop item 1/80th of the time
	x := rand.Intn(80)
	if x != 0 {
		log.Println("no item drop this turn")
		return nil
	}

	is, err := models.AllItems()
	if err != nil {
		return err
	}

	// Get random item to drop
	i := rand.Intn(len(is))
	item := is[i]

	// Get random player to find it
	ps, err := models.AllPlayers()
	if err != nil {
		return err
	}
	i = rand.Intn(len(ps))
	p := ps[i]

	// Player finds dropped item
	err = p.PickUpItem(item)
	if err != nil {
		return err
	}

	frontend.PushMessage(droppedItemMessage(p, item))
	return nil
}

func droppedItemMessage(p models.Player, item models.Item) string {
	var itemAndArticle string
	if len(item.Article) > 0 {
		itemAndArticle = fmt.Sprintf("%v %v", item.Article, item.Name)
	} else {
		itemAndArticle = item.Name
	}

	var incrs []string
	if item.LevelIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v levels", item.LevelIncr))
	}
	if item.XPIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v xp", item.XPIncr))
	}
	if item.StrengthIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v strength", item.StrengthIncr))
	}
	if item.CharismaIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v charisma", item.CharismaIncr))
	}
	if item.IntellectIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v intellect", item.IntellectIncr))
	}
	if item.AgilityIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v agility", item.AgilityIncr))
	}
	if item.LuckIncr > 0 {
		incrs = append(incrs, fmt.Sprintf("%v luck", item.LuckIncr))
	}

	if len(incrs) > 1 {
		lastIndex := len(incrs) - 1
		incrs[lastIndex] = "and " + incrs[lastIndex]
	}

	msg := fmt.Sprintf(
		"%v found %v! They gain %v ",
		p.NameAndTitle(), itemAndArticle,
		strings.Join(incrs, ", "),
	)

	msg += fmt.Sprintf(
		"\n%v's new stats: "+
			"*Level*: %v *XP*: %v *Strength*: %v *Charisma*: %v "+
			"*Intellect*: %v *Agility*: %v *Luck*: %v",
		k.NameAndTitle(),
		p.Level,
		p.XP,
		p.Strength,
		p.Charisma,
		p.Intellect,
		p.Agility,
		p.Luck,
	)

	return msg
}
