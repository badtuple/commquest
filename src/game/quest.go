package game

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"../db"
	"../frontends"
	"../quest"
)

func maybeStartOrProgressQuest() error {
	q := quest.InProgressOrGen()

	// If the q.ID is 0 then it's new. Otherwise it's already inprogress
	if q.ID == 0 {
		err := q.Create()
		if err != nil {
			return err
		}

		log.Println("starting a new quest")
	} else {
		log.Printf("resuming quest %v", q.ID)
	}

	if !q.Intro.Completed {
		msg := q.Intro.RenderMsg()
		err := frontend.PushMessage(msg)
		if err != nil {
			log.Printf("could not push message for quest intro: %v", err.Error())
			return err
		}

		q.Intro.Completed = true
		return q.Save()
	}

	if len(q.Journey) > 0 {
		for i, o := range q.Journey {
			if o.Completed {
				continue
			}

			msg := o.RenderMsg()
			err := frontend.PushMessage(msg)
			if err != nil {
				log.Printf("could not push message for quest obstacle: %v", err.Error())
				return err
			}

			o.Completed = true
			q.Journey[i] = o

			return q.Save()
		}
	}

	if q.Return != nil && !q.Return.Completed {
		msg := q.Return.RenderMsg()
		err := frontend.PushMessage(msg)
		if err != nil {
			return err
		}

		q.Return.Completed = true
		return q.Save()
	}

	if !q.End.Completed {
		msg := q.End.RenderMsg()
		err := frontend.PushMessage(msg)
		q.End.Completed = true

		err = rewardPlayers(q)
		if err != nil {
			return err
		}

		q.InProgress = false
		return q.Save()
	}

	return errors.New(fmt.Sprintf("unaccounted for quest state: %v", q.ID))
}

func rewardPlayers(q quest.Quest) error {
	var in []string
	for _, id := range q.PlayerIDs {
		in = append(in, strconv.Itoa(id))
	}

	rewardXP := q.RewardXP()
	inStr := strings.Join(in, ",")
	query := fmt.Sprintf(`UPDATE players SET xp = xp + $1 WHERE id IN (%v)`, inStr)
	_, err := db.PSQL().Exec(query, rewardXP)
	if err != nil {
		return err
	}

	ps, err := q.Players()
	if err != nil {
		return err
	}

	var pnames []string
	for _, p := range ps {
		pnames = append(pnames, p.Handle)
	}

	msg := fmt.Sprintf(`%v completed the quest and each got %v xp!`, strings.Join(pnames, ", "), rewardXP)
	return frontend.PushMessage(msg)
}
