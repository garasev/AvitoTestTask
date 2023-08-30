package repository

import "github.com/garasev/AvitoTestTask/internal/models"

type Repository interface {
	AddSlug(slug models.Slug) (int, error)
	GetSlug(id int) (models.Slug, error)
}
