package parser

import (
	"fmt"
	"service/internal/pkg/model"
	"strconv"
	"time"
)

var dayNames = map[time.Weekday]string{
	time.Sunday:    "Minggu",
	time.Monday:    "Senin",
	time.Tuesday:   "Selasa",
	time.Wednesday: "Rabu",
	time.Thursday:  "Kamis",
	time.Friday:    "Jumat",
	time.Saturday:  "Sabtu",
}

var monthNames = map[time.Month]string{
	time.January:   "Januari",
	time.February:  "Februari",
	time.March:     "Maret",
	time.April:     "April",
	time.May:       "Mei",
	time.June:      "Juni",
	time.July:      "Juli",
	time.August:    "Agustus",
	time.September: "September",
	time.October:   "Oktober",
	time.November:  "November",
	time.December:  "Desember",
}

func FormatReportDate(t time.Time) string {
	return fmt.Sprintf("%s, %d %s %d",
		dayNames[t.Weekday()],
		t.Day(),
		monthNames[t.Month()],
		t.Year(),
	)
}

func FormatChannelMessage(user model.ReportUser, report model.Report) string {
	completedYesterday := ""
	if report.CompletedYesterday != "" {
		completedYesterday, _ = strconv.Unquote(report.CompletedYesterday)
	}

	planToday := ""
	if report.PlanToday != "" {
		planToday, _ = strconv.Unquote(report.PlanToday)
	}

	finishEstimation := ""
	if report.FinishEstimation != "" {
		finishEstimation, _ = strconv.Unquote(report.FinishEstimation)
	}

	blockers := ""
	if report.Blockers != "" {
		blockers, _ = strconv.Unquote(report.Blockers)
	}

	mood := ""
	if report.Mood != "" {
		mood, _ = strconv.Unquote(report.Mood)
	}

	return fmt.Sprintf(
		"📋 *Daily Report — %s — %s*\n\n"+
			"✅ *What did you complete yesterday?*\n%s\n\n"+
			"🔨 *What will you do today?*\n%s\n\n"+
			"⏰ *When will you be finished?*\n%s\n\n"+
			"🚧 *Anything blocking your progress?*\n%s\n\n"+
			"😊 *How do you feel today?*\n%s",
		user.Name,
		FormatReportDate(report.ReportDate),
		completedYesterday,
		planToday,
		finishEstimation,
		blockers,
		mood,
	)
}
