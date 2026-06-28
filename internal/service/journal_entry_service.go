package service

import (
	"bank/internal/models"
	"bank/internal/repository"
	"time"
)

type JournalEntryService struct {
	repo *repository.JournalEntryRepo
}

func NewJournalEntryService(repo *repository.JournalEntryRepo) *JournalEntryService {
	return &JournalEntryService{repo: repo}
}

func (s *JournalEntryService) FindAll() ([]models.JournalEntry, error) {
	return s.repo.FindAll()
}

func (s *JournalEntryService) FindByDate(date time.Time) ([]models.JournalEntry, error) {
	return s.repo.FindByDate(date)
}

func (s *JournalEntryService) Create(e *models.JournalEntry) error {
	return s.repo.Create(e)
}

func (s *JournalEntryService) Delete(id string) error {
	return s.repo.Delete(id)
}
