package config

import (
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/socketmode"
	"os"
)

var (
	SlackAPI    *slack.Client
	SlackSocket *socketmode.Client
)

func InitSlack() {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")

	SlackAPI = slack.New(
		botToken,
		slack.OptionAppLevelToken(appToken),
	)

	SlackSocket = socketmode.New(SlackAPI)
}
