package parser

import (
	"service/internal/pkg/model"
	"time"
)

// ── Response structs ──────────────────────────────────────────────────────────

// MetaResponse holds pagination metadata returned in every list response.
type MetaResponse struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}

// UserSummary is a compact user representation embedded in report responses.
type UserSummary struct {
	ID      uint   `json:"id"`
	SlackID string `json:"slackId"`
	Name    string `json:"name"`
	Email   string `json:"email"`
}

// ReportResponse is the standard single-report representation for list and
// detail endpoints.
type ReportResponse struct {
	ID                 uint        `json:"id"`
	User               UserSummary `json:"user"`
	ReportDate         string      `json:"reportDate"`
	CompletedYesterday string      `json:"completedYesterday"`
	PlanToday          string      `json:"planToday"`
	FinishEstimation   string      `json:"finishEstimation"`
	Blockers           string      `json:"blockers"`
	Mood               string      `json:"mood"`
	CreatedAt          time.Time   `json:"createdAt"`
	CompletedAt        string      `json:"completedAt"`
}

// UserReportsResponse groups a user with all their reports (used by
// GET /reports/by-user).
type UserReportsResponse struct {
	User    UserSummary      `json:"user"`
	Reports []ReportResponse `json:"reports"`
	Count   int              `json:"count"`
}

// DateReportsResponse groups reports that share the same report date (used by
// GET /reports/by-date).
type DateReportsResponse struct {
	Date    string           `json:"date"`
	Reports []ReportResponse `json:"reports"`
	Count   int              `json:"count"`
}

// ── Conversion helpers ────────────────────────────────────────────────────────

func toUserSummary(u model.ReportUser) UserSummary {
	return UserSummary{
		ID:      u.ID,
		SlackID: u.SlackID,
		Name:    u.Name,
		Email:   u.Email,
	}
}

// ToReportResponse converts a single model.Report to ReportResponse.
// The Report's User association must be preloaded.
func ToReportResponse(r model.Report) ReportResponse {
	completedYesterday := ""
	if r.CompletedYesterday != "" {
		completedYesterday = r.CompletedYesterday
		//completedYesterday, _ = strconv.Unquote(r.CompletedYesterday)
	}

	planToday := ""
	if r.PlanToday != "" {
		planToday = r.PlanToday
		//planToday, _ = strconv.Unquote(r.PlanToday)
	}

	finishEstimation := ""
	if r.FinishEstimation != "" {
		finishEstimation = r.FinishEstimation
		//finishEstimation, _ = strconv.Unquote(r.FinishEstimation)
	}

	blockers := ""
	if r.Blockers != "" {
		blockers = r.Blockers
		//blockers, _ = strconv.Unquote(r.Blockers)
	}

	mood := ""
	if r.Mood != "" {
		mood = r.Mood
		//mood, _ = strconv.Unquote(r.Mood)
	}

	var completedAt string
	if r.SessionCompletedAt != nil {
		completedAt = r.SessionCompletedAt.Format("2006-01-02 15:04:05")
	}

	return ReportResponse{
		ID:                 r.ID,
		User:               toUserSummary(r.User),
		ReportDate:         r.ReportDate.Format("2006-01-02"),
		CompletedYesterday: completedYesterday,
		PlanToday:          planToday,
		FinishEstimation:   finishEstimation,
		Blockers:           blockers,
		Mood:               mood,
		CreatedAt:          r.CreatedAt.In(time.Local),
		CompletedAt:        completedAt,
	}
}

// ToReportResponses converts a slice of model.report.
func ToReportResponses(reports []model.Report) []ReportResponse {
	out := make([]ReportResponse, len(reports))
	for i, r := range reports {
		out[i] = ToReportResponse(r)
	}
	return out
}

// GroupByUser groups a flat slice of reports into per-user buckets,
// preserving the existing order within each user's reports.
func GroupByUser(reports []model.Report) []UserReportsResponse {
	// Use an ordered slice of keys to maintain insertion order.
	type bucket struct {
		user    UserSummary
		reports []ReportResponse
	}

	indexMap := make(map[uint]int) // userID → index in buckets
	var buckets []bucket

	for _, r := range reports {
		rr := ToReportResponse(r)
		idx, exists := indexMap[r.UserID]
		if !exists {
			idx = len(buckets)
			indexMap[r.UserID] = idx
			buckets = append(buckets, bucket{user: toUserSummary(r.User)})
		}
		buckets[idx].reports = append(buckets[idx].reports, rr)
	}

	out := make([]UserReportsResponse, len(buckets))
	for i, b := range buckets {
		out[i] = UserReportsResponse{
			User:    b.user,
			Reports: b.reports,
			Count:   len(b.reports),
		}
	}
	return out
}

// GroupByDate groups a flat slice of reports into per-date buckets,
// preserving the existing order within each date's reports.
func GroupByDate(reports []model.Report) []DateReportsResponse {
	type bucket struct {
		date    string
		reports []ReportResponse
	}

	indexMap := make(map[string]int)
	var buckets []bucket

	for _, r := range reports {
		rr := ToReportResponse(r)
		dateKey := r.ReportDate.Format("2006-01-02")
		idx, exists := indexMap[dateKey]
		if !exists {
			idx = len(buckets)
			indexMap[dateKey] = idx
			buckets = append(buckets, bucket{date: dateKey})
		}
		buckets[idx].reports = append(buckets[idx].reports, rr)
	}

	out := make([]DateReportsResponse, len(buckets))
	for i, b := range buckets {
		out[i] = DateReportsResponse{
			Date:    b.date,
			Reports: b.reports,
			Count:   len(b.reports),
		}
	}
	return out
}
