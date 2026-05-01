package reportservice

import (
	"errors"
	"fmt"
	xtremepkg "github.com/globalxtreme/go-core/v2/pkg"
	"log"
	"os"
	"time"

	"service/internal/pkg/constant"
	apperror "service/internal/pkg/error"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	"service/internal/pkg/thirdparty"
	reportrepository "service/internal/report/repository"
)

type ReportService struct {
	repo      *reportrepository.ReportRepository
	slack     *thirdparty.Slack
	channelID string
}

func New(repo *reportrepository.ReportRepository, slack *thirdparty.Slack, channelID string) *ReportService {
	return &ReportService{repo: repo, slack: slack, channelID: channelID}
}

// ProcessMessage handles an incoming DM from a Slack user.
// Returns the bot's reply message.
func (s *ReportService) ProcessMessage(user *model.ReportUser, text string) (string, error) {
	today := todayDate()

	session, err := s.repo.FindSessionByUserAndDate(user.ID, today)

	if err != nil && !errors.Is(err, apperror.ErrSessionNotFound) {
		return "", err
	}

	// Session selesai → tolak pengisian ulang
	if session != nil && session.IsCompleted {
		return "Kamu sudah mengisi daily report hari ini ✅", nil
	}

	// Tidak ada session → jawaban pertama langsung jadi step 1
	if session == nil {
		return s.startSessionWithAnswer(user, today, text)
	}

	// Session ada dan belum selesai → simpan jawaban dan lanjut
	return s.processAnswer(user, session, text, today)
}

func (s *ReportService) startSessionWithAnswer(user *model.ReportUser, today time.Time, answer string) (string, error) {
	session := &model.ReportSession{
		UserID:      user.ID,
		ReportDate:  today,
		CurrentStep: constant.SLACK_STEP_1,
	}
	if err := s.repo.CreateSession(session); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	report := &model.Report{
		UserID:     user.ID,
		ReportDate: today,
	}
	if err := s.repo.CreateReport(report); err != nil {
		return "", fmt.Errorf("create report: %w", err)
	}

	// Simpan jawaban step 1 langsung dari pesan pertama user
	applyAnswer(report, constant.SLACK_STEP_1, answer)
	if err := s.repo.SaveReport(report); err != nil {
		return "", fmt.Errorf("save step 1: %w", err)
	}

	// Advance session ke step 2
	session.CurrentStep = constant.SLACK_STEP_2
	if err := s.repo.SaveSession(session); err != nil {
		return "", fmt.Errorf("advance to step 2: %w", err)
	}

	return constant.SlackStepQuestions[constant.SLACK_STEP_2], nil
}

func (s *ReportService) processAnswer(user *model.ReportUser, session *model.ReportSession, answer string, today time.Time) (string, error) {
	report, err := s.repo.FindByUserAndDate(user.ID, today)
	if err != nil {
		return "", err
	}

	currentStep := int(session.CurrentStep)
	applyAnswer(report, currentStep, answer)

	if err := s.repo.SaveReport(report); err != nil {
		return "", fmt.Errorf("save report: %w", err)
	}

	if currentStep == constant.SLACK_STEP_5 {
		return s.completeSession(user, session, report)
	}

	session.CurrentStep++
	if err := s.repo.SaveSession(session); err != nil {
		return "", fmt.Errorf("advance session: %w", err)
	}

	return constant.SlackStepQuestions[int(session.CurrentStep)], nil
}

func (s *ReportService) completeSession(user *model.ReportUser, session *model.ReportSession, report *model.Report) (string, error) {
	now := time.Now()
	xtremepkg.LogInfo(now)
	session.IsCompleted = true
	session.CompletedAt = &now

	if err := s.repo.SaveSession(session); err != nil {
		return "", fmt.Errorf("complete session: %w", err)
	}

	// Kirim ke channel Slack di background
	go s.sendToChannel(user, report)

	return fmt.Sprintf("✅ Daily report kamu sudah tersimpan. \nKamu dapat melihat report pada dashboard: %s. \nTerima kasih!", os.Getenv("FRONTEND_HOST")), nil
}

func (s *ReportService) sendToChannel(user *model.ReportUser, report *model.Report) {
	msg := parser.FormatChannelMessage(*user, *report)

	ts, err := s.slack.SendChannelMessage(s.channelID, msg)
	if err != nil {
		log.Printf("[ReportService] sendToChannel error (user=%s): %v", user.SlackID, err)
		return
	}

	report.SlackChannelTs = ts
	if err := s.repo.SaveReport(report); err != nil {
		log.Printf("[ReportService] save slackChannelTs error: %v", err)
	}
}

func applyAnswer(report *model.Report, step int, answer string) {
	switch step {
	case constant.SLACK_STEP_1:
		report.CompletedYesterday = answer
	case constant.SLACK_STEP_2:
		report.PlanToday = answer
	case constant.SLACK_STEP_3:
		report.FinishEstimation = answer
	case constant.SLACK_STEP_4:
		report.Blockers = answer
	case constant.SLACK_STEP_5:
		report.Mood = answer
	}
}

func todayDate() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}
