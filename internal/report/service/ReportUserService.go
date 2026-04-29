package reportservice

import (
	"service/internal/pkg/model"
	reportrepository "service/internal/report/repository"
)

type ReportUserService struct {
	repo *reportrepository.ReportUserRepository
}

func NewUser(repo *reportrepository.ReportUserRepository) *ReportUserService {
	return &ReportUserService{repo: repo}
}

func (s *ReportUserService) GetBySlackID(slackID string) (*model.ReportUser, error) {
	return s.repo.FindBySlackID(slackID)
}

func (s *ReportUserService) GetAllActive() ([]model.ReportUser, error) {
	return s.repo.FindAllActive()
}

func (s *ReportUserService) Create(user *model.ReportUser) error {
	return s.repo.Create(user)
}
