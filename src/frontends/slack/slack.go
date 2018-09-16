package slack_frontend

import (
	"database/sql"
	"fmt"

	"../../config"
	"../../models"
	"../../util"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus"
)

var log *logrus.Entry = util.LoggerFor("frnt")

type SlackFrontend struct {
	Client *slack.Client
	RTM    *slack.RTM
}

func (_ SlackFrontend) Name() string {
	return "slack"
}

func (fe SlackFrontend) PushMessage(msg string) error {
	cid := config.Get().Slack.ChannelID
	_, _, err := fe.RTM.PostMessage(cid, msg, slack.NewPostMessageParameters())
	return err
}

func (fe SlackFrontend) Serve() error {
	key := config.Get().Slack.APIKey
	fe.Client = slack.New(key)
	fe.RTM = fe.Client.NewRTM()

	go fe.RTM.ManageConnection()

Loop:
	for {
		select {
		case msg := <-fe.RTM.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				log.Println("connected to slack")

			case *slack.ConnectedEvent:
				log.Printf("connected %+v times", ev.ConnectionCount)

			case *slack.RTMError:
				log.Printf("rtm error: %s", ev.Error())

			case *slack.InvalidAuthEvent:
				log.Print("invalid credentials")
				break Loop

			case *slack.PresenceChangeEvent:
				//log.Printf("presence_changed_event: %+v", ev)

			case *slack.MessageEvent:
				log.Printf("Message: %v", ev)

			case *slack.MemberJoinedChannelEvent:
				log.Printf("member_joined_channel_event: %+v", ev)
				fe.memberJoinedChannelHandler(ev)

			case *slack.MemberLeftChannelEvent:
				// TODO: impl
				//log.Printf("member_left_channel_event: %+v", ev)

			default:
				log.Printf("unhandled event type: %+v", msg.Type)
			}
		}
	}

	return nil
}

func (fe SlackFrontend) memberJoinedChannelHandler(ev *slack.MemberJoinedChannelEvent) {
	uid := ev.User
	log.Printf("member_joined_channel: %+v", uid)

	user, err := fe.Client.GetUserInfo(uid)
	if err != nil {
		log.Printf("could not get user %v: %v", uid, err.Error())
		return
	}

	log.Printf("%+v", user)
	handle := getHandleFromSlackUser(user)
	p, err := models.FindPlayerByHandle(handle)
	if err == nil {
		log.Printf("user by handle %v already exists", handle)
		return
	}

	if err != sql.ErrNoRows {
		log.Println("%+v", err.Error())
		return
	}

	title := user.Profile.Title
	p, err = models.CreatePlayer(handle, handle, title)
	if err != nil {
		log.Println("could not create player: %v", err.Error())
		return
	}

	msg := createdPlayerMessage(p)
	err = fe.PushMessage(msg)
	if err != nil {
		log.Printf("could not send message to channel: %v", err.Error())
		return
	}

	log.Printf("created player: %+v", p.Handle)
}

// Slack profiles can be filled out to varying degrees
// so we fallback to less ideal but more likely to be
// filled out values
func getHandleFromSlackUser(user *slack.User) string {
	display := user.Profile.DisplayName
	if len(display) > 0 {
		return display
	}

	realName := user.Profile.RealName
	if len(realName) > 0 {
		return realName
	}

	return user.Name
}

func createdPlayerMessage(p models.Player) string {
	name := p.Handle
	if len(p.Class) > 0 {
		name += " the " + p.Class
	}
	return fmt.Sprintf(`%v has joined the game!`, name)
}
