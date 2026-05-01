package reportservice

import (
	"fmt"
	"service/internal/pkg/docx"
	"service/internal/pkg/form"
	"service/internal/pkg/model"
	"service/internal/pkg/parser"
	reportrepository "service/internal/report/repository"
)

type ReportExportService struct {
	repo *reportrepository.ReportQueryRepository
}

func NewExport(repo *reportrepository.ReportQueryRepository) *ReportExportService {
	return &ReportExportService{repo: repo}
}

// Export builds and returns a .docx file as []byte.
// The groupBy field in the form controls the document structure:
//   - "none"  → flat list, one report per section
//   - "user"  → grouped by employee name
//   - "date"  → grouped by report date
func (s *ReportExportService) Export(f *form.ReportExportForm) ([]byte, error) {
	reports, err := s.repo.FindAllForExport(f.FromDate, f.ToDate, f.Sort)
	if err != nil {
		return nil, fmt.Errorf("export: fetch reports: %w", err)
	}

	b := docx.New()
	b.AddTitle("Daily Report Export")

	// Subtitle: date range
	subtitle := "All dates"
	if f.FromDate != "" || f.ToDate != "" {
		from := f.FromDate
		if from == "" {
			from = "—"
		}
		to := f.ToDate
		if to == "" {
			to = "—"
		}
		subtitle = fmt.Sprintf("%s  s/d  %s", from, to)
	}
	b.AddCentered(subtitle)
	b.AddBlank()

	switch f.GroupBy {
	case "user":
		buildGroupedByUser(b, reports)
	case "date":
		buildGroupedByDate(b, reports)
	default:
		buildFlat(b, reports)
	}

	return b.Build()
}

// ── document builders ─────────────────────────────────────────────────────────

func buildFlat(b *docx.Builder, reports []model.Report) {
	for i, r := range reports {
		if i > 0 {
			b.AddSeparator()
			b.AddBlank()
		}
		appendReportBlock(b, r)
	}
}

func buildGroupedByUser(b *docx.Builder, reports []model.Report) {
	// Group preserving order.
	type bucket struct {
		userName string
		reports  []model.Report
	}

	indexMap := make(map[uint]int)
	var buckets []bucket

	for _, r := range reports {
		idx, ok := indexMap[r.UserID]
		if !ok {
			idx = len(buckets)
			indexMap[r.UserID] = idx
			buckets = append(buckets, bucket{userName: r.User.Name})
		}
		buckets[idx].reports = append(buckets[idx].reports, r)
	}

	for i, bkt := range buckets {
		if i > 0 {
			b.AddBlank()
		}
		b.AddBold(fmt.Sprintf("👤 %s", bkt.userName))
		b.AddSeparator()
		for j, r := range bkt.reports {
			if j > 0 {
				b.AddBlank()
			}
			appendReportBlock(b, r)
		}
	}
}

func buildGroupedByDate(b *docx.Builder, reports []model.Report) {
	type bucket struct {
		date    string
		reports []model.Report
	}

	indexMap := make(map[string]int)
	var buckets []bucket

	for _, r := range reports {
		dateKey := r.ReportDate.Format("2006-01-02")
		idx, ok := indexMap[dateKey]
		if !ok {
			idx = len(buckets)
			indexMap[dateKey] = idx
			buckets = append(buckets, bucket{date: dateKey})
		}
		buckets[idx].reports = append(buckets[idx].reports, r)
	}

	for i, bkt := range buckets {
		if i > 0 {
			b.AddBlank()
		}
		b.AddBold(fmt.Sprintf("📅 %s", parser.FormatReportDate(buckets[i].reports[0].ReportDate)))
		b.AddSeparator()
		for j, r := range bkt.reports {
			if j > 0 {
				b.AddBlank()
			}
			appendReportBlock(b, r)
		}
	}
}

// appendReportBlock writes one report's fields into the document.
func appendReportBlock(b *docx.Builder, r model.Report) {
	b.AddNormal(fmt.Sprintf("%s — %s", r.User.Name, r.ReportDate.Format("2006-01-02")))
	b.AddLabelValue("What did you complete yesterday?", r.CompletedYesterday)
	b.AddLabelValue("What will you do today?", r.PlanToday)
	b.AddLabelValue("When will you be finished?", r.FinishEstimation)
	b.AddLabelValue("Anything blocking your progress?", r.Blockers)
	b.AddLabelValue("How do you feel today?", r.Mood)
}
