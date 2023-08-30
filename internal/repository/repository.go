package repository

import "github.com/garasev/AvitoTestTask/internal/models"

type Repository interface {
	//AddArchive(slug models.Slug) error
	AddSlug(slug models.Slug) (int, error)
	GetSlug(id int) (models.Slug, error)
	GetSlugs() ([]models.Slug, error)
}
