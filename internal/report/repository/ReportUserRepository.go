package reportrepository

import (
	"errors"

	apperror "service/internal/pkg/error"
	"service/internal/pkg/model"

	"gorm.io/gorm"
)

type ReportUserRepository struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) *ReportUserRepository {
	return &ReportUserRepository{db: db}
}

func (r *ReportUserRepository) FindBySlackID(slackID string) (*model.ReportUser, error) {
	var user model.ReportUser
	err := r.db.Where(`"slackId" = ?`, slackID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *ReportUserRepository) FindAllActive() ([]model.ReportUser, error) {
	var users []model.ReportUser
	err := r.db.Where(`"isActive" = ?`, true).Find(&users).Error
	return users, err
}

func (r *ReportUserRepository) FindByID(id uint) (*model.ReportUser, error) {
	var user model.ReportUser
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *ReportUserRepository) Create(user *model.ReportUser) error {
	return r.db.Create(user).Error
}

func (r *ReportUserRepository) Save(user *model.ReportUser) error {
	return r.db.Save(user).Error
}
