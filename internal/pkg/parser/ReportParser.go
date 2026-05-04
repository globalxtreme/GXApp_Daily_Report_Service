package parser

import (
	"fmt"
	"service/internal/pkg/model"
	"time"
)

var dayNames = map[time.Weekday]string{
	time.Sunday:    "Sunday",
	time.Monday:    "Monday",
	time.Tuesday:   "Tuesday",
	time.Wednesday: "Wednesday",
	time.Thursday:  "Thursday",
	time.Friday:    "Friday",
	time.Saturday:  "Saturday",
}

var monthNames = map[time.Month]string{
	time.January:   "January",
	time.February:  "February",
	time.March:     "March",
	time.April:     "April",
	time.May:       "May",
	time.June:      "June",
	time.July:      "July",
	time.August:    "August",
	time.September: "September",
	time.October:   "October",
	time.November:  "November",
	time.December:  "December",
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
	return fmt.Sprintf(
		"📋 *Daily Report — %s — %s*\n\n"+
			"✅ *What did you complete yesterday?*\n%s\n\n"+
			"🔨 *What will you do today?*\n%s\n\n"+
			"⏰ *When will you be finished?*\n%s\n\n"+
			"🚧 *Anything blocking your progress?*\n%s\n\n"+
			"😊 *How do you feel today?*\n%s",
		user.Name,
		FormatReportDate(report.ReportDate),
		report.CompletedYesterday,
		report.PlanToday,
		report.FinishEstimation,
		report.Blockers,
		report.Mood,
	)
}
