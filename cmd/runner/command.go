package runner

import (
	"service/internal/app/console"
)

func init() {
	console.RegisterCommand(rootCmd)
}
