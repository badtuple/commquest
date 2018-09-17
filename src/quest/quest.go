package quest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"strings"

	"../db"
	"../models"
)

// A Quest's narrative consists of 4 sections:
//	Intro -- a Call to action for the group.
//		Examples:
//			a villager needs an item
//			a town is being terrorized
//			a person needs rescued
//
//	Journey -- an array of obstacles that need to be overcome.
//		Each array element takes 1 turn to complete.
//		Examples:
//			a dungeon must be explored
//			a person must be defeated
//			a person must be found
//
//	Return -- Completely optional step where an additional
//		Journey element must be passed through on the
//		return to the origin of the quest.
//
//	End -- The original Call to action is confronted and
//		the characters are rewarded with XP, items, etc.
//		Examples:
//			villager accepts needed item
//			town is grateful that you've saved them
//			person is safely returned home
type Quest struct {
	ID         int    `db:"id"`    // DB row id
	Raw        string `db:"state"` // Serialized raw state
	InProgress bool   `db:"inprogress"`

	// Narrative structure
	Intro   intro      `json:"intro"`
	Journey []obstacle `json:"journey"`
	Return  *obstacle  `json:"return"`
	End     end        `json:"end"`

	PlayerIDs []int `json:"player_ids"`
}

// Serialize the current state of the quest into JSON.
// Each turn we pull out the quest and make the next move.
// We then either mark the Quest as complete or reserialize
// so that it can be used during the next turn.

func (q *Quest) Create() error {
	raw, err := json.Marshal(q)
	if err != nil {
		return err
	}
	q.Raw = string(raw)

	err = db.PSQL().QueryRowx(`
		INSERT INTO quests (state, inprogress)
		VALUES $1, true
		RETURNING id, state, inprogress
	`, q.Raw).StructScan(&q)
	if err != nil {
		return err
	}

	return json.Unmarshal(raw, &q)
}

func (q *Quest) Save() error {
	raw, err := json.Marshal(q)
	if err != nil {
		return err
	}
	q.Raw = string(raw)

	_, err = db.PSQL().Exec(`
		UPDATE quests
		SET state = $1, inprogress = $2
		WHERE id = $3
	`, q.Raw, q.InProgress, q.ID)
	return err
}

func (q Quest) Players() ([]models.Player, error) {
	var in []string
	for _, id := range q.PlayerIDs {
		in = append(in, strconv.Itoa(id))
	}

	inStr := strings.Join(in, ",")
	query := fmt.Sprintf(`SELECT 
		SELECT id, handle, name, class, xp, level,
			strength, charisma, intellect, agility, luck,
			created_at, updated_at
		FROM players
		WHERE id IN (%v)`, inStr)

	var ps []models.Player
	err := db.PSQL().Select(&ps, query)
	return ps, err
}

func (q Quest) RewardXP() int {
	// super simple for now...we just count up number of
	// turns and give players 100 per turn.  Definitely
	// a TODO to make this more involved.

	var turns int = 2 // intro and end
	turns += len(q.Journey)
	if q.Return != nil {
		turns += 1
	}

	return turns * 100
}

type intro struct {
	// Possible intro types:
	//	obtain    -- Obtain an item or animal
	//	destroy   -- Destroy an item
	//	transport -- Transport an item
	//	save_town -- Save a town that's being terroized
	//	rescue    -- Rescue a person or animal
	//	escort    -- Escort a person/animal between towns
	Type string `json:"type"`

	Tmpl string `json:"tmpl"` // Intro text template

	// Item to be obtained, destroyed, or transported
	Item *models.Item `json:"item"`

	// Person to be rescued or escorted
	Person *models.Person `json:"person"`

	// Town to be saved. Alternatively the town where
	// the quest originated.
	Town *string `json:"town"`

	// Certain intro types require a specific ending
	EndType []string `json:"end_type"`

	Completed bool `json:"completed"`
}

func (in intro) RenderMsg() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(in.Tmpl)
	if err != nil {
		panic(err.Error())
	}

	err = tmpl.Execute(buf, in)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

type obstacle struct {
	Completed bool   `json:"completed"`
	Tmpl      string `json:"tmpl"`
}

func (o obstacle) RenderMsg() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(o.Tmpl)
	if err != nil {
		panic(err.Error())
	}

	err = tmpl.Execute(buf, o)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

type end struct {
	// Possible end types:
	//	none   -- Quest just ends.
	// 	thanks -- The person who sent you on the quest says thank you.
	//	party  -- The town celebrates their safety.
	Type string `json:"type"`
	Tmpl string `json:"tmpl"` // Ending text template

	Completed bool `json:"completed"`
}

func (e end) RenderMsg() string {
	buf := new(bytes.Buffer)
	tmpl, err := template.New("").Parse(e.Tmpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(buf, e)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// TODOs:
//	[ ] obtain -- Obtain a particular item or animal
//	[ ] destroy -- Destroy a particular item
//	[ ] transport -- Transport a particular item
//	[ ] save_town -- Save a town that's being terroized
//	[ ] rescue -- Rescue a particular person or animal
//	[ ] escort -- Transport a person or animal between towns

func Gen() Quest {
	return genObtainQuest()
}

func InProgressOrGen() Quest {
	var q Quest
	err := db.PSQL().
		QueryRowx(`
		SELECT id, state, inprogress
		FROM quests
		WHERE inprogress = true
		LIMIT 1
		`).
		StructScan(&q)
	if err == sql.ErrNoRows {
		return Gen()
	}

	if err != nil {
		log.Printf("could not get inprogress quest: %v", err.Error())
		return Gen()
	}

	return q
}
