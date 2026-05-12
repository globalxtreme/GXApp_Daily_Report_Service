package command

import (
	"fmt"
	"log"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"

	"service/internal/pkg/config"
	"service/internal/pkg/constant"
	"service/internal/pkg/thirdparty"
	reportrepository "service/internal/report/repository"
)

type SendReminderCommand struct{}

func (c *SendReminderCommand) Command(cobraCmd *cobra.Command) {
	var slackId string

	addCommand := cobra.Command{
		Use:  "send-reminder",
		Long: "Manually send a daily report reminder to a specific user by their Slack ID",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			DBClose := config.InitDB()
			defer DBClose()

			config.InitSlack()

			c.handle(slackId)
		},
	}

	addCommand.Flags().StringVar(&slackId, "slackId", "", "Slack ID of the user to send reminder to (required)")
	_ = addCommand.MarkFlagRequired("slackId")

	cobraCmd.AddCommand(&addCommand)
}

func (c *SendReminderCommand) Handle() {}

func (c *SendReminderCommand) handle(slackId string) {
	repo := reportrepository.NewUser(config.PgSQL)

	user, err := repo.FindBySlackID(slackId)
	if err != nil {
		log.Fatalf("[SendReminder] User dengan slackId '%s' tidak ditemukan: %v", slackId, err)
	}

	if !user.IsActive {
		log.Fatalf("[SendReminder] User '%s' (%s) tidak aktif, reminder tidak dikirim", user.Name, user.SlackID)
	}

	slackClient := thirdparty.New(config.SlackAPI)

	msg := fmt.Sprintf(
		"Hey %s! It's time to fill in your daily report for today. 📋\n\n*%s*",
		user.Name,
		constant.SlackStepQuestions[constant.SLACK_STEP_1],
	)

	if err := slackClient.SendDM(user.SlackID, msg); err != nil {
		log.Fatalf("[SendReminder] Gagal mengirim DM ke %s (%s): %v", user.Name, user.SlackID, err)
	}

	fmt.Printf("Reminder berhasil dikirim ke %s (%s)\n", user.Name, user.SlackID)
}
