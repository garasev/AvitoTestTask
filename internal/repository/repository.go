package repository

import (
	"time"

	"github.com/garasev/AvitoTestTask/internal/models"
)

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

	AddSlugsUser(id int, slugs []models.Slug, duration time.Duration) error
	DeleteSlugsUser(id int, slugs []models.Slug) error

	GetUserSlugs(id int) ([]models.Slug, error)

	GetCntRandomUsers(cntRandomUsers int) ([]int, error)
	GetCntUsers() (int, error)

	GetUserArchive(id int, date time.Time) ([]models.Archive, error)
}
