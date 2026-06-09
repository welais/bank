package service

import (
	"bank/internal/models"
	"bank/internal/repository"
)

type AccountService struct {
	repo *repository.AccountRepo
}

func NewAccountService(repo *repository.AccountRepo) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) Open(clientID string, accountType string) (*models.Account, error) {
	return s.repo.Open(clientID, accountType)
}

func (s *AccountService) Close(id string) error {
	return s.repo.Close(id)
}

func (s *AccountService) FindAll() ([]models.Account, error) {
	return s.repo.FindAll()
}

func (s *AccountService) FindByNumber(number string) (*models.Account, error) {
	return s.repo.FindByNumber(number)
}

func (s *AccountService) FindByClientLastName(lastName string) ([]models.Account, error) {
	return s.repo.FindByClientLastName(lastName)
}

func (s *AccountService) DeleteAllClosed() (int64, error) {
	return s.repo.DeleteAllClosed()
}

func (s *AccountService) Delete(id string) error {
	return s.repo.Delete(id)
}
