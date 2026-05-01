package reportrepository

import (
	"errors"
	apperror "service/internal/pkg/error"
	"service/internal/pkg/model"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ReportRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

// --- Session ---

func (r *ReportRepository) FindSessionByUserAndDate(userID uint, date time.Time) (*model.ReportSession, error) {
	var session model.ReportSession
	err := r.db.Where(`"userId" = ? AND "reportDate" = ?`, userID, date).First(&session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrSessionNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *ReportRepository) CreateSession(session *model.ReportSession) error {
	return r.db.Create(session).Error
}

func (r *ReportRepository) SaveSession(session *model.ReportSession) error {
	return r.db.Debug().Save(session).Error
}

// --- Report ---

func (r *ReportRepository) FindByUserAndDate(userID uint, date time.Time) (*model.Report, error) {
	var report model.Report
	err := r.db.Where(`"userId" = ? AND "reportDate" = ?`, userID, date).First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrReportNotFound
		}
		return nil, err
	}
	return &report, nil
}

func (r *ReportRepository) CreateReport(report *model.Report) error {
	return r.db.Clauses(clause.OnConflict{DoNothing: true}).Create(report).Error
}

func (r *ReportRepository) SaveReport(report *model.Report) error {
	return r.db.Save(report).Error
}
