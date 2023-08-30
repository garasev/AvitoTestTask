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

func (s *Service) AddSlug(slug models.Slug) (int, error) {
	return s.repo.AddSlug(slug)
}

func (s *Service) DeleteSlug(slug models.Slug) error {
	return s.repo.DeleteSlug(slug)
}
