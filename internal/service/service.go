package service

import (
	"github.com/garasev/AvitoTestTask/internal/models"
	"github.com/garasev/AvitoTestTask/internal/repository"
)

type Service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetSlug(id int) (models.Slug, error) {
	return s.repo.GetSlug(id)
}

func (s *Service) GetSlugs() ([]models.Slug, error) {
	return s.repo.GetSlugs()
}

func (s *Service) AddSlug(slug models.Slug) error {
	return s.repo.AddSlug(slug)
}

func (s *Service) DeleteSlug(slug models.Slug) error {
	return s.repo.DeleteSlug(slug)
}

func (s *Service) GetUsers() ([]models.AvitoUser, error) {
	return s.repo.GetUsers()
}

func (s *Service) AddUsers(cnt int) error {
	return s.repo.AddUsers(cnt)
}

func (s *Service) CheckUser(id int) (bool, error) {
	return s.repo.CheckUserExist(id)
}

func (s *Service) CheckSlugs(slugs []models.Slug) (bool, error) {
	return s.repo.CheckSlugsExist(slugs)
}

func (s *Service) AddSlugsUser(id int, slugs []models.Slug) error {
	return s.repo.AddSlugsUser(id, slugs)
}

func (s *Service) DeleteSlugsUser(id int, slugs []models.Slug) error {
	return s.repo.DeleteSlugsUser(id, slugs)
}

func (s *Service) GetUserSlugs(id int) ([]models.Slug, error) {
	return s.repo.GetUserSlugs(id)
}
