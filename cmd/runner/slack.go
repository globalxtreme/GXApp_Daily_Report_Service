package runner

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"github.com/go-co-op/gocron"
	"github.com/spf13/cobra"

	slackapp "service/internal/app/slack"
	"service/internal/pkg/config"
	"service/internal/pkg/thirdparty"
	reportrepository "service/internal/report/repository"
	reportservice "service/internal/report/service"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:  "xtreme:slack",
		Long: "Running Slack Bot (Socket Mode) + cron reminder scheduler",
		Run: func(cmd *cobra.Command, args []string) {
			xtremepkg.InitDevMode()
			config.InitTZ()

			DBClose := config.InitDB()
			defer DBClose()

			config.InitSlack()

			// Wire up dependencies
			slackWrapper := thirdparty.New(config.SlackAPI)

			userRepo := reportrepository.NewUser(config.PgSQL)
			reportRepo := reportrepository.New(config.PgSQL)

			userSvc := reportservice.NewUser(userRepo)
			reportSvc := reportservice.New(reportRepo, slackWrapper, os.Getenv("SLACK_CHANNEL_ID"))

			botHandler := slackapp.New(config.SlackSocket, slackWrapper, userSvc, reportSvc)

			// Start gocron scheduler for daily reminders
			reminderTime := os.Getenv("REMINDER_TIME")
			if reminderTime == "" {
				reminderTime = "14:00"
			}

			sch := gocron.NewScheduler(time.Local)
			_, err := sch.Every(1).Week().
				Monday().Tuesday().Wednesday().Thursday().Friday().
				At(reminderTime).
				Do(func() {
					log.Println("[Scheduler] Sending daily reminders...")
					botHandler.SendReminders(userSvc)
				})
			if err != nil {
				log.Fatalf("[Scheduler] Failed to register reminder job: %v", err)
			}

			sch.StartAsync()
			log.Printf("[Scheduler] Reminders scheduled at %s (Mon-Fri)", reminderTime)
			defer sch.Stop()

			// Start Slack Socket Mode (blocking) in a goroutine
			// then wait for OS signal for graceful shutdown
			go botHandler.Run()

			quit := make(chan os.Signal, 1)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
			<-quit

			log.Println("[SlackBot] Shutting down...")
		},
	})
}
