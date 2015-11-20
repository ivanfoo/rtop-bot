package bot

import (
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/daneharrigan/hipchat"
	"github.com/ivanfoo/rtop-bot/utils"
	"github.com/nlopes/slack"
)

type Bot struct {
	botOptions BotOptions
	SystemUser *user.User
}

type BotOptions struct {
	Username   string
	SSHKeyPath string
	SlackToken string
}

func NewBot(bopts BotOptions) *Bot {
	b := new(Bot)
	b.botOptions = bopts

	if bopts.Username != "" {
		b.SystemUser, _ = user.Lookup(bopts.Username)
	} else {
		b.SystemUser, _ = user.Current()
	}

	return b
}

func (b *Bot) DoSlack() {
	api := slack.New(b.botOptions.SlackToken)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	mention := ""
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			mention = "<@" + ev.Info.User.ID + ">"
			if ev.ConnectionCount == 1 {
				log.Printf("bot [%s] ready", ev.Info.User.Name)
				log.Print("hit ^C to exit")
			} else {
				log.Printf("bot [%s] reconnected", ev.Info.User.Name)
			}
		case *slack.MessageEvent:
			if strings.HasPrefix(ev.Msg.Text, mention) {
				t := strings.TrimPrefix(ev.Msg.Text, mention)
				go func(text, ch string) {
					r := b.process(text)
					rtm.SendMessage(rtm.NewOutgoingMessage(r, ch))
				}(t, ev.Msg.Channel)
			}
		case *slack.InvalidAuthEvent:
			log.Print("bad Slack API token")
			os.Exit(1)
		}
	}
}

func getUserInfo(client *hipchat.Client, id string) (string, string) {
	id = id + "@chat.hipchat.com"
	client.RequestUsers()
	select {
	case users := <-client.Users():
		for _, user := range users {
			if user.Id == id {
				log.Printf("using username [%s] and mention name [%s]",
					user.Name, user.MentionName)
				return user.Name, user.MentionName
			}
		}
	case <-time.After(10 * time.Second):
		log.Print("timed out waiting for user list")
		os.Exit(1)
	}
	return "rtop-bot", "rtop-bot"
}

func (b *Bot) process(request string) string {

	parts := strings.Fields(request)
	if len(parts) != 3 || parts[1] != "status" {
		return "say \"status <hostname>\" to see vital stats of <hostname>"
	}

	hostname := utils.CleanHostname(parts[2])
	return utils.SSHConnectTmp(b.botOptions.Username, hostname, b.botOptions.SSHKeyPath)
	/*
		parts[2] = utils.CleanHostname(parts[2])
		address, username, keypath := getSshEntryOrDefault(parts[2])
		client, err := sshConnect(username, address, keypath, userHome)
		if err != nil {
			return fmt.Sprintf("[%s]: %v", parts[2], err)
		}

		stats := Stats{}
		getAllStats(client, &stats)
		result := fmt.Sprintf(
			`[%s] up %s, load %s %s %s, procs %s running of %s total
				[%s] mem: %s of %s free, swap %s of %s free
				`,
			stats.Hostname, fmtUptime(&stats), stats.Load1, stats.Load5,
			stats.Load10, stats.RunningProcs, stats.TotalProcs,
			stats.Hostname, fmtBytes(stats.MemFree), fmtBytes(stats.MemTotal),
			fmtBytes(stats.SwapFree), fmtBytes(stats.SwapTotal),
		)
		if len(stats.FSInfos) > 0 {
			for _, fs := range stats.FSInfos {
				result += fmt.Sprintf("[%s] fs %s: %s of %s free\n",
					stats.Hostname,
					fs.MountPoint,
					fmtBytes(fs.Free),
					fmtBytes(fs.Used+fs.Free),
				)
			}
		}
		return result
	*/
}
