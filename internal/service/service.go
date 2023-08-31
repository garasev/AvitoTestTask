package service

import (
	"fmt"
	"time"

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

func (s *Service) AddSlug(slug models.Slug, percent int) ([]int, error) {
	var users []int
	err := s.repo.AddSlug(slug)
	if err != nil {
		return nil, err
	}

	cnt, err := s.repo.GetCntUsers()
	if err != nil {
		return nil, err
	}
	fmt.Println(cnt)
	cntRandomUser := (cnt * percent) / 100
	users, err = s.repo.GetCntRandomUsers(cntRandomUser)
	if err != nil {
		return nil, err
	}
	fmt.Println(cntRandomUser)
	fmt.Println(users)
	slugs := []models.Slug{slug}
	for _, userId := range users {
		err = s.AddSlugsUser(userId, slugs, 0)
		if err != nil {
			return users, err
		}
	}

	return users, nil
}

func (s *Service) DeleteSlug(slug models.Slug) error {
	userSlugs, err := s.repo.GetUserBySlug(slug)
	if err != nil {
		return err
	}
	for _, userSlug := range userSlugs {
		err := s.repo.AddArchive(userSlug.UserId, []models.Slug{slug}, false)
		if err != nil {
			return err
		}
	}
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

func (s *Service) AddSlugsUser(id int, slugs []models.Slug, duration time.Duration) error {
	err := s.repo.AddArchive(id, slugs, true)
	if err != nil {
		return err
	}
	return s.repo.AddSlugsUser(id, slugs, duration)
}

func (s *Service) DeleteSlugsUser(id int, slugs []models.Slug) error {
	err := s.repo.AddArchive(id, slugs, false)
	if err != nil {
		return err
	}
	return s.repo.DeleteSlugsUser(id, slugs)
}

func (s *Service) GetUserSlugs(id int) ([]models.Slug, error) {
	return s.repo.GetUserSlugs(id)
}

func (s *Service) GetUserArchive(id int, date time.Time) ([]models.Archive, error) {
	return s.repo.GetUserArchive(id, date)
}
