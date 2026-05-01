package reportrepository

import (
	"errors"

	apperror "service/internal/pkg/error"
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

// sessionSelect adalah SELECT clause yang menyertakan completedAt dari report_sessions
// sebagai alias sessionCompletedAt agar dipetakan ke model.Report.SessionCompletedAt.
const sessionSelect = `reports.*, report_sessions."completedAt" as "sessionCompletedAt"`

// sessionJoin adalah JOIN clause untuk menggabungkan report_sessions.
const sessionJoin = `JOIN report_sessions ON report_sessions."userId" = reports."userId" AND report_sessions."reportDate" = reports."reportDate"`

type ReportQueryRepository struct {
	db *gorm.DB
}

func NewQuery(db *gorm.DB) *ReportQueryRepository {
	return &ReportQueryRepository{db: db}
}

// completedBase returns a scoped query that joins report_sessions
// and filters only completed reports.
func (r *ReportQueryRepository) completedBase() *gorm.DB {
	return r.db.Model(&model.Report{}).
		Joins(sessionJoin).
		Where(`report_sessions."isCompleted" = ?`, true)
}

// applyDateRange appends optional date range conditions on reports."reportDate".
func applyDateRange(q *gorm.DB, fromDate, toDate string) *gorm.DB {
	if fromDate != "" {
		q = q.Where(`reports."reportDate" >= ?`, fromDate)
	}
	if toDate != "" {
		q = q.Where(`reports."reportDate" <= ?`, toDate)
	}
	return q
}

// ── FindAll ─────────────────────────────────────────────────────────────────

// FindAll returns a paginated list of completed reports sorted by createdAt.
func (r *ReportQueryRepository) FindAll(page, limit int, sort string) ([]model.Report, int64, error) {
	var total int64
	if err := r.completedBase().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var reports []model.Report
	err := r.completedBase().
		Preload("User").
		Select(sessionSelect).
		Order(`reports."createdAt" ` + sort).
		Offset((page-1)*limit).
		Limit(limit).
		Find(&reports).Error

	return reports, total, err
}

// ── FindByUserGrouped ────────────────────────────────────────────────────────

// FindByUserGrouped returns completed reports within a date range.
// Pagination is over distinct users; total reflects the number of users.
// Reports within each user are sorted by sortBy / sort.
func (r *ReportQueryRepository) FindByUserGrouped(
	fromDate, toDate, sortBy, sort string,
	page, limit int,
) ([]model.Report, int64, error) {

	// Count distinct users that have completed reports in the date range.
	var total int64
	if err := applyDateRange(r.completedBase(), fromDate, toDate).
		Distinct(`reports."userId"`).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch a page of distinct user IDs.
	var userIDs []uint
	if err := applyDateRange(r.completedBase(), fromDate, toDate).
		Select(`reports."userId"`).
		Group(`reports."userId"`).
		Offset((page - 1) * limit).
		Limit(limit).
		Pluck(`reports."userId"`, &userIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(userIDs) == 0 {
		return []model.Report{}, total, nil
	}

	// Choose the order column inside each user's group.
	orderCol := `"createdAt"`
	if sortBy == "reportDate" {
		orderCol = `"reportDate"`
	}

	var reports []model.Report
	err := applyDateRange(
		r.db.Model(&model.Report{}).
			Preload("User").
			Joins(sessionJoin).
			Where(`reports."userId" IN ? AND report_sessions."isCompleted" = ?`, userIDs, true).
			Select(sessionSelect),
		fromDate, toDate,
	).Order(`reports."userId", reports.` + orderCol + ` ` + sort).
		Find(&reports).Error

	return reports, total, err
}

// ── FindByDateGrouped ────────────────────────────────────────────────────────

// FindByDateGrouped returns completed reports within a date range.
// Pagination is over distinct report dates; total reflects the number of dates.
func (r *ReportQueryRepository) FindByDateGrouped(
	fromDate, toDate string,
	page, limit int,
) ([]model.Report, int64, error) {

	// Count distinct report dates.
	var total int64
	if err := applyDateRange(r.completedBase(), fromDate, toDate).
		Distinct(`reports."reportDate"`).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch a page of distinct dates (as strings for simplicity).
	var dates []string
	if err := applyDateRange(r.completedBase(), fromDate, toDate).
		Select(`reports."reportDate"::text`).
		Group(`reports."reportDate"`).
		Order(`reports."reportDate" DESC`).
		Offset((page - 1) * limit).
		Limit(limit).
		Pluck(`reports."reportDate"::text`, &dates).Error; err != nil {
		return nil, 0, err
	}

	if len(dates) == 0 {
		return []model.Report{}, total, nil
	}

	var reports []model.Report
	err := r.db.Model(&model.Report{}).
		Preload("User").
		Joins(sessionJoin).
		Where(`reports."reportDate"::text IN ? AND report_sessions."isCompleted" = ?`, dates, true).
		Select(sessionSelect).
		Order(`reports."reportDate" DESC, reports."userId" ASC`).
		Find(&reports).Error

	return reports, total, err
}

// ── FindByID ─────────────────────────────────────────────────────────────────

// FindByID returns a single completed report by its primary key.
func (r *ReportQueryRepository) FindByID(id uint) (*model.Report, error) {
	var report model.Report
	err := r.completedBase().
		Preload("User").
		Select(sessionSelect).
		Where(`reports."id" = ?`, id).
		First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrReportNotFound
		}
		return nil, err
	}
	return &report, nil
}

// ── FindAllForExport ─────────────────────────────────────────────────────────

// FindAllForExport returns all completed reports for export (no pagination).
// fromDate / toDate are optional; sort is "asc" | "desc" applied to reportDate.
func (r *ReportQueryRepository) FindAllForExport(fromDate, toDate, sort string) ([]model.Report, error) {
	q := applyDateRange(r.completedBase(), fromDate, toDate).
		Preload("User").
		Select(sessionSelect).
		Order(`reports."reportDate" ` + sort + `, reports."userId" ASC`)

	var reports []model.Report
	err := q.Find(&reports).Error
	return reports, err
}
