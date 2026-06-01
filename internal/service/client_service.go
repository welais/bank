package service

import (
	"bank/internal/models"
	"bank/internal/repository"
)

type ClientService struct {
	repo *repository.ClientRepo
}

func NewClientService(r *repository.ClientRepo) *ClientService {
	return &ClientService{repo: r}
}

func (s *ClientService) Create(c *models.Client) error {
	return s.repo.Create(c)
}

func (s *ClientService) FindAll() ([]models.Client, error) {
	return s.repo.FindAll()
}

func (s *ClientService) FindByLastName(lastName string) ([]models.Client, error) {
	return s.repo.FindByLastName(lastName)
}

func (s *ClientService) FindByPhone(phone string) ([]models.Client, error) {
	return s.repo.FindByPhone(phone)
}

func (s *ClientService) Update(c *models.Client) error {
	return s.repo.Update(c)
}

func (s *ClientService) Delete(id string) error {
	return s.repo.Delete(id)
}
