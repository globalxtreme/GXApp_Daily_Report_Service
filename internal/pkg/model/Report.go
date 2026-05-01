package model

import (
	xtrememodel "github.com/globalxtreme/go-core/v2/model"
	"time"
)

type Report struct {
	xtrememodel.BaseModel
	UserID             uint       `gorm:"not null;column:userId;uniqueIndex:idx_reports_user_date"`
	User               ReportUser `gorm:"foreignKey:UserID"`
	ReportDate         time.Time  `gorm:"type:date;not null;column:reportDate;uniqueIndex:idx_reports_user_date"`
	CompletedYesterday string     `gorm:"type:text;column:completedYesterday"`
	PlanToday          string     `gorm:"type:text;column:planToday"`
	FinishEstimation   string     `gorm:"type:text;column:finishEstimation"`
	Blockers           string     `gorm:"type:text;column:blockers"`
	Mood               string     `gorm:"type:text;column:mood"`
	SlackChannelTs     string     `gorm:"size:50;column:slackChannelTs"`
}

func (Report) TableName() string {
	return "reports"
}
