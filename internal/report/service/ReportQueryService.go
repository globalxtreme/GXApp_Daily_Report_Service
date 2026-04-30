package reportservice

import (
	"service/internal/pkg/form"
	"service/internal/pkg/parser"
	reportrepository "service/internal/report/repository"
)

type ReportQueryService struct {
	repo *reportrepository.ReportQueryRepository
}

func NewQuery(repo *reportrepository.ReportQueryRepository) *ReportQueryService {
	return &ReportQueryService{repo: repo}
}

// List returns a paginated flat list of completed reports.
func (s *ReportQueryService) List(f *form.ReportListForm) ([]parser.ReportResponse, parser.MetaResponse, error) {
	reports, total, err := s.repo.FindAll(f.Page, f.Limit, f.Sort)
	if err != nil {
		return nil, parser.MetaResponse{}, err
	}

	return parser.ToReportResponses(reports), parser.MetaResponse{
		Page:  f.Page,
		Limit: f.Limit,
		Total: total,
	}, nil
}

// ListByUser returns completed reports grouped by user within a date range.
func (s *ReportQueryService) ListByUser(f *form.ReportByUserForm) ([]parser.UserReportsResponse, parser.MetaResponse, error) {
	reports, total, err := s.repo.FindByUserGrouped(f.FromDate, f.ToDate, f.SortBy, f.Sort, f.Page, f.Limit)
	if err != nil {
		return nil, parser.MetaResponse{}, err
	}

	return parser.GroupByUser(reports), parser.MetaResponse{
		Page:  f.Page,
		Limit: f.Limit,
		Total: total,
	}, nil
}

// ListByDate returns completed reports grouped by date within a date range.
func (s *ReportQueryService) ListByDate(f *form.ReportByDateForm) ([]parser.DateReportsResponse, parser.MetaResponse, error) {
	reports, total, err := s.repo.FindByDateGrouped(f.FromDate, f.ToDate, f.Page, f.Limit)
	if err != nil {
		return nil, parser.MetaResponse{}, err
	}

	return parser.GroupByDate(reports), parser.MetaResponse{
		Page:  f.Page,
		Limit: f.Limit,
		Total: total,
	}, nil
}

// GetByID returns a single completed report by ID.
func (s *ReportQueryService) GetByID(id uint) (*parser.ReportResponse, error) {
	report, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	resp := parser.ToReportResponse(*report)
	return &resp, nil
}
