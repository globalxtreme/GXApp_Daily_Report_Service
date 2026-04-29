package model

import "time"

type ReportUser struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	SlackID   string    `gorm:"uniqueIndex;not null;size:50;column:slackId"`
	Name      string    `gorm:"not null;type:text;column:name"`
	Email     string    `gorm:"type:text;column:email"`
	IsActive  bool      `gorm:"default:true;column:isActive"`
	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
}

func (ReportUser) TableName() string {
	return "report_users"
}
