package command

import (
	"fmt"
	"log"
	"service/internal/pkg/constant"
	"service/internal/pkg/model"
	"service/internal/pkg/thirdparty"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"

	"service/internal/pkg/config"
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
	//_ = addCommand.MarkFlagRequired("slackId")

	cobraCmd.AddCommand(&addCommand)
}

func (c *SendReminderCommand) Handle() {}

func (c *SendReminderCommand) handle(slackId string) {
	repo := reportrepository.NewUser(config.PgSQL)

	users := make([]model.ReportUser, 0)
	if slackId != "" {
		user, err := repo.FindBySlackID(slackId)
		if err != nil {
			log.Fatalf("[SendReminder] User with slackId '%s' not found: %v", slackId, err)
		}

		if !user.IsActive {
			log.Fatalf("[SendReminder] User '%s' (%s) is not active, Please contact manager", user.Name, user.SlackID)
		}
	} else {
		users, _ = repo.FindAllActive()
	}

	if len(users) > 0 {
		slackClient := thirdparty.New(config.SlackAPI)

		for _, user := range users {
			msg := fmt.Sprintf(
				"Hey %s! It's time to fill in your daily report for today. 📋\n\n*%s*",
				user.Name,
				constant.SlackStepQuestions[constant.SLACK_STEP_1],
			)

			if err := slackClient.SendDM(user.SlackID, msg); err != nil {
				log.Fatalf("[SendReminder] Failed to send DM to %s (%s): %v", user.Name, user.SlackID, err)
			}

			fmt.Printf("Reminder successfully sent to %s (%s)\n", user.Name, user.SlackID)
		}
	}
}
