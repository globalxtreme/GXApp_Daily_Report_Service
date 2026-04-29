package config

import (
	"log"
	"os"
	"time"
)

func InitTZ() {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "Asia/Makassar"
	}

	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Panicf("Error loading location: %v\n", err)
	}

	time.Local = loc
}
