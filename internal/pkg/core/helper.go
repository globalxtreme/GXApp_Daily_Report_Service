package core

import (
	"math/rand"
	"net/url"
	"os"
	"strings"
	"time"
)

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "production"
}

func SetDateRangeFromAPI(parameters url.Values) (time.Time, time.Time) {
	return SetDateRange(parameters.Get("fromDate"), parameters.Get("toDate"))
}

func SetDateRange(fromDateArgs string, toDateArgs string) (time.Time, time.Time) {
	var fromDate, toDate time.Time

	now := time.Now()

	if len(fromDateArgs) > 0 {
		fromDate, _ = time.Parse("02/01/2006 15:04:05", fromDateArgs+" 00:00:00")
	} else {
		fromDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, -1, 0)
	}

	if len(toDateArgs) > 0 {
		toDate, _ = time.Parse("02/01/2006 15:04:05", toDateArgs+" 23:59:59")
	} else {
		toDate = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
	}

	return fromDate, toDate
}

func StrPadLeft(original string, padLength int, padChar rune) string {
	if len(original) >= padLength {
		return original
	}
	padding := strings.Repeat(string(padChar), padLength-len(original))
	return padding + original
}

func RandInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
