package service

import (
	"bank/internal/models"
	"bank/internal/repository"
	"time"
)

type StatementService struct {
	repo *repository.StatementRepo
}

func NewStatementService(repo *repository.StatementRepo) *StatementService {
	return &StatementService{repo: repo}
}

func (s *StatementService) GetByAccount(accountID string, from, to time.Time) ([]models.AccountStatement, error) {
	return s.repo.GetByAccount(accountID, from, to)
}

func (s *StatementService) GetAll(from, to time.Time) ([]models.AccountStatement, error) {
	return s.repo.GetAll(from, to)
}
