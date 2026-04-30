package form

// ReportListForm adalah query params untuk GET /reports.
type ReportListForm struct {
	Sort  string `query:"sort"`  // "asc" | "desc", default: "desc"
	Page  int    `query:"page"`  // default: 1
	Limit int    `query:"limit"` // default: 20
}

func (f *ReportListForm) Defaults() {
	if f.Sort != "asc" && f.Sort != "desc" {
		f.Sort = "desc"
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
}

func (f *ReportListForm) Offset() int {
	return (f.Page - 1) * f.Limit
}

// ReportByUserForm adalah query params untuk GET /reports/by-user.
type ReportByUserForm struct {
	FromDate string `query:"fromDate"` // "YYYY-MM-DD", required
	ToDate   string `query:"toDate"`   // "YYYY-MM-DD", required
	SortBy   string `query:"sortBy"`   // "createdAt" | "reportDate", default: "createdAt"
	Sort     string `query:"sort"`     // "asc" | "desc", default: "desc"
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

func (f *ReportByUserForm) Defaults() {
	if f.SortBy != "reportDate" {
		f.SortBy = "createdAt"
	}
	if f.Sort != "asc" && f.Sort != "desc" {
		f.Sort = "desc"
	}
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
}

func (f *ReportByUserForm) Offset() int {
	return (f.Page - 1) * f.Limit
}

// ReportByDateForm adalah query params untuk GET /reports/by-date.
type ReportByDateForm struct {
	FromDate string `query:"fromDate"` // "YYYY-MM-DD", required
	ToDate   string `query:"toDate"`   // "YYYY-MM-DD", required
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

func (f *ReportByDateForm) Defaults() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.Limit < 1 {
		f.Limit = 20
	}
}

func (f *ReportByDateForm) Offset() int {
	return (f.Page - 1) * f.Limit
}

// ReportExportForm adalah query params untuk GET /reports/export.
type ReportExportForm struct {
	FromDate string `query:"fromDate"` // "YYYY-MM-DD", optional
	ToDate   string `query:"toDate"`   // "YYYY-MM-DD", optional
	Sort     string `query:"sort"`     // "asc" | "desc", default: "desc"
	GroupBy  string `query:"groupBy"`  // "none" | "user" | "date", default: "none"
}

func (f *ReportExportForm) Defaults() {
	if f.Sort != "asc" && f.Sort != "desc" {
		f.Sort = "desc"
	}
	if f.GroupBy != "user" && f.GroupBy != "date" {
		f.GroupBy = "none"
	}
}
