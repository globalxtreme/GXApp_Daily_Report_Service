package form

type ReportUserForm struct {
	SlackID string `json:"slackId" validate:"required"`
	Name    string `json:"name"    validate:"required"`
	Email   string `json:"email"`
}

type ReportFilterForm struct {
	UserID     uint   `json:"userId"`
	ReportDate string `json:"reportDate"`
}
