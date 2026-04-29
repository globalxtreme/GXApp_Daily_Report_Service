package model

import "time"

type ReportSession struct {
	ID          uint       `gorm:"primaryKey;column:id"`
	UserID      uint       `gorm:"not null;column:userId;uniqueIndex:idx_sessions_user_date"`
	User        ReportUser `gorm:"foreignKey:UserID"`
	ReportDate  time.Time  `gorm:"type:date;not null;column:reportDate;uniqueIndex:idx_sessions_user_date"`
	CurrentStep int16      `gorm:"default:1;column:currentStep"`
	IsCompleted bool       `gorm:"default:false;column:isCompleted"`
	StartedAt   time.Time  `gorm:"column:startedAt;autoCreateTime"`
	CompletedAt *time.Time `gorm:"column:completedAt"`
}

func (ReportSession) TableName() string {
	return "report_sessions"
}
