package error

import "errors"

var (
	ErrUserNotFound        = errors.New("user not found")
	ErrUserNotActive       = errors.New("user is not active")
	ErrSessionAlreadyDone  = errors.New("daily report sudah diisi hari ini")
	ErrSessionNotFound     = errors.New("session not found")
	ErrReportNotFound      = errors.New("report not found")
)
