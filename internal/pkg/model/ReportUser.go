package model

import xtrememodel "github.com/globalxtreme/go-core/v2/model"

type ReportUser struct {
	xtrememodel.BaseModel
	SlackID  string `gorm:"uniqueIndex;not null;size:50;column:slackId"`
	Name     string `gorm:"not null;type:text;column:name"`
	Email    string `gorm:"type:text;column:email"`
	IsActive bool   `gorm:"default:true;column:isActive"`
}

func (ReportUser) TableName() string {
	return "report_users"
}
