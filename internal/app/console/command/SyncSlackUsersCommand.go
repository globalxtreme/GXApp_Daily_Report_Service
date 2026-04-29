package command

import (
	"errors"
	"fmt"
	"log"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/spf13/cobra"

	"service/internal/pkg/config"
	apperror "service/internal/pkg/error"
	"service/internal/pkg/model"
	reportrepository "service/internal/report/repository"
)

type SyncSlackUsersCommand struct{}

func (c *SyncSlackUsersCommand) Command(cobraCmd *cobra.Command) {
	addCommand := cobra.Command{
		Use:  "sync-slack-users",
		Long: "Sync active Slack workspace members to report_users table",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()

			config.InitTZ()

			DBClose := config.InitDB()
			defer DBClose()

			config.InitSlack()

			c.Handle()
		},
	}

	cobraCmd.AddCommand(&addCommand)
}

func (c *SyncSlackUsersCommand) Handle() {
	// 1. Ambil semua user dari Slack workspace
	slackUsers, err := config.SlackAPI.GetUsers()
	if err != nil {
		log.Fatalf("[SyncSlackUsers] Slack API error: %v", err)
	}

	repo := reportrepository.NewUser(config.PgSQL)

	added := 0
	skipped := 0

	for _, u := range slackUsers {
		// 2. Filter: skip bot, deleted, dan akun slackbot
		if u.IsBot || u.Deleted || u.Name == "slackbot" {
			continue
		}

		// 3. Cek apakah slackId sudah ada di DB
		_, err := repo.FindBySlackID(u.ID)
		if err == nil {
			// Sudah ada → skip
			skipped++
			continue
		}
		if !errors.Is(err, apperror.ErrUserNotFound) {
			// Error DB yang tidak terduga
			log.Printf("[SyncSlackUsers] DB lookup error (slackId=%s): %v", u.ID, err)
			continue
		}

		// Belum ada → tentukan nama terbaik
		name := u.Profile.RealName
		if name == "" {
			name = u.Name
		}

		newUser := &model.ReportUser{
			SlackID:  u.ID,
			Name:     name,
			Email:    u.Profile.Email,
			IsActive: true,
		}

		if err := repo.Create(newUser); err != nil {
			log.Printf("[SyncSlackUsers] Insert error (slackId=%s, name=%s): %v", u.ID, name, err)
			continue
		}

		fmt.Printf("  + Added: %-30s  slack_id=%s\n", name, u.ID)
		added++
	}

	// 4. Summary
	fmt.Printf("\nSync selesai: %d user ditambahkan, %d user sudah ada sebelumnya\n", added, skipped)
}
