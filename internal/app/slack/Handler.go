package slack

import (
	"errors"
	"fmt"
	slacklib "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"github.com/slack-go/slack/socketmode"
	"log"

	"service/internal/pkg/constant"
	apperror "service/internal/pkg/error"
	"service/internal/pkg/thirdparty"
	reportservice "service/internal/report/service"
)

type Handler struct {
	socket        *socketmode.Client
	slackClient   *thirdparty.Slack
	userService   *reportservice.ReportUserService
	reportService *reportservice.ReportService
}

func New(
	socket *socketmode.Client,
	slackClient *thirdparty.Slack,
	userService *reportservice.ReportUserService,
	reportService *reportservice.ReportService,
) *Handler {
	return &Handler{
		socket:        socket,
		slackClient:   slackClient,
		userService:   userService,
		reportService: reportService,
	}
}

// Run starts the Socket Mode event loop. This call blocks until the connection closes.
func (h *Handler) Run() {
	go h.listenEvents()

	log.Println("[SlackBot] Connected, listening for events...")
	if err := h.socket.Run(); err != nil {
		log.Fatalf("[SlackBot] Socket run error: %v", err)
	}
}

func (h *Handler) listenEvents() {
	for evt := range h.socket.Events {
		switch evt.Type {
		case socketmode.EventTypeConnecting:
			log.Println("[SlackBot] Connecting to Slack...")

		case socketmode.EventTypeConnected:
			log.Println("[SlackBot] Connected to Slack via Socket Mode!")

		case socketmode.EventTypeEventsAPI:
			apiEvent, ok := evt.Data.(slackevents.EventsAPIEvent)
			if !ok {
				continue
			}
			h.socket.Ack(*evt.Request)
			h.handleAPIEvent(apiEvent)

		case socketmode.EventTypeErrorBadMessage:
			log.Printf("[SlackBot] Bad message error: %v", evt.Data)
		}
	}
}

func (h *Handler) handleAPIEvent(event slackevents.EventsAPIEvent) {
	if event.Type != slackevents.CallbackEvent {
		return
	}

	switch ev := event.InnerEvent.Data.(type) {
	case *slackevents.MessageEvent:
		// Hanya proses DM dari user nyata (bukan bot, bukan system message)
		if ev.BotID == "" && ev.SubType == "" && ev.ChannelType == "im" {
			h.handleDM(ev)
		}
	}
}

func (h *Handler) handleDM(ev *slackevents.MessageEvent) {
	user, err := h.userService.GetBySlackID(ev.User)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			_ = h.slackClient.SendDM(ev.User,
				"Kamu belum terdaftar sebagai pengguna daily report. "+
					"Hubungi admin untuk mendaftarkan akun Slack kamu.",
			)
			return
		}
		log.Printf("[SlackBot] GetBySlackID error (slackId=%s): %v", ev.User, err)
		return
	}

	if !user.IsActive {
		_ = h.slackClient.SendDM(ev.User, "Akunmu sedang tidak aktif. Hubungi admin.")
		return
	}

	reply, err := h.reportService.ProcessMessage(user, ev.Text)
	if err != nil {
		log.Printf("[SlackBot] ProcessMessage error (userId=%d): %v", user.ID, err)
		_ = h.slackClient.SendDM(ev.User,
			fmt.Sprintf("Terjadi kesalahan: %v\nSilakan coba lagi.", err),
		)
		return
	}

	if err := h.slackClient.SendDM(ev.User, reply); err != nil {
		log.Printf("[SlackBot] SendDM error (slackId=%s): %v", ev.User, err)
	}
}

// SendReminders mengirim pesan reminder ke semua user aktif.
// Dipanggil oleh scheduler.
func (h *Handler) SendReminders(userSvc *reportservice.ReportUserService) {
	users, err := userSvc.GetAllActive()
	if err != nil {
		log.Printf("[SlackBot] SendReminders: fetch users error: %v", err)
		return
	}

	sent := 0
	for _, user := range users {
		// Reminder langsung embed pertanyaan step 1 agar user bisa langsung menjawab
		msg := fmt.Sprintf(
			"Halo %s! Waktunya mengisi daily report kamu hari ini 📋\n\n*%s*",
			user.Name,
			constant.SlackStepQuestions[constant.SLACK_STEP_1],
		)
		if err := h.slackClient.SendDM(user.SlackID, msg); err != nil {
			log.Printf("[SlackBot] SendReminders: failed to DM %s: %v", user.SlackID, err)
			continue
		}
		sent++
	}

	log.Printf("[SlackBot] Reminders sent: %d/%d", sent, len(users))
}

// API returns the underlying Slack API client (untuk kebutuhan lain di luar paket ini).
func (h *Handler) API() *slacklib.Client {
	return &h.socket.Client
}
