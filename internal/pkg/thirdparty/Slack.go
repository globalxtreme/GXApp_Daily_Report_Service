package thirdparty

import "github.com/slack-go/slack"

type Slack struct {
	api *slack.Client
}

func New(api *slack.Client) *Slack {
	return &Slack{api: api}
}

func (c *Slack) SendDM(userSlackID, message string) error {
	_, _, err := c.api.PostMessage(
		userSlackID,
		slack.MsgOptionText(message, false),
	)
	return err
}

func (c *Slack) SendChannelMessage(channelID, message string) (string, error) {
	_, ts, err := c.api.PostMessage(
		channelID,
		slack.MsgOptionText(message, false),
	)
	return ts, err
}
