package repository

import "github.com/garasev/AvitoTestTask/internal/models"

type Repository interface {
	//AddArchive(slug models.Slug) error
	AddSlug(slug models.Slug) error
	GetSlug(id int) (models.Slug, error)
	GetSlugs() ([]models.Slug, error)
	DeleteSlug(slug models.Slug) error

	AddUsers(cnt int) error
	GetUsers() ([]models.AvitoUser, error)

	CheckUserExist(id int) (bool, error)
	CheckSlugsExist([]models.Slug) (bool, error)

	AddSlugsUser(id int, slugs []models.Slug) error
	DeleteSlugsUser(id int, slugs []models.Slug) error
}
