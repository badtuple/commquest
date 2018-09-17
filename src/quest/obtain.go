package quest

import (
	"log"
	"math/rand"

	"../models"
)

func genObtainQuest() Quest {
	q := Quest{
		Intro:   genObtainIntro(),
		Journey: nil,
		Return:  nil,
	}

	q.End = genObtainEnd(q)
	return q
}

func genObtainIntro() intro {
	in := intro{
		Type: "obtain",
		Tmpl: `{{.Person.Name}} stops {{.Players}} and asks them to find {{.Item.Article}} {{.Item.Name}}`,

		Item:   nil,
		Person: nil,
		Town:   nil,
		EndType: []string{
			"none",
			"thanks",
			"party",
		},
	}

	p := models.GenPerson()
	in.Person = &p

	item, err := models.GetRandomItem()
	in.Item = item
	if err != nil {
		log.Printf("could not get random item: %v", err.Error())
		in.Item = &models.Item{
			ID:      0,
			Name:    "artifact",
			Article: "an",
		}
	}

	return in
}

func genObtainEnd(q Quest) end {
	var typ string = "none"
	if len(q.Intro.EndType) > 0 {
		i := rand.Intn(len(q.Intro.EndType))
		typ = q.Intro.EndType[i]
	}

	var tmpl string
	switch typ {
	case "none":
		tmpl = ``
	case "thanks":
		tmpl = `{{.Person.Name}} accepts the {{.Item.Name}} and is extremely thankful!`
	case "party":
		tmpl = `{{.Person.Name}} is overjoyed and throws a celebration in your honor.`
	}

	return end{
		Type: typ,
		Tmpl: tmpl,
	}
}
