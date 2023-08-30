package postgresql

import (
	"database/sql"
	"strconv"

	"github.com/garasev/AvitoTestTask/internal/models"
)

type PostgresqlRep struct {
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB) *PostgresqlRep {
	return &PostgresqlRep{
		DB: conn,
	}
}

func (r *PostgresqlRep) AddSlug(slug models.Slug) (int, error) {
	var id int
	querySql := `INSERT INTO slug (name) VALUES ($1) RETURNING id`

	err := r.DB.QueryRow(
		querySql,
		slug.Name,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *PostgresqlRep) GetSlug(id int) (models.Slug, error) {
	var slug models.Slug
	querySql := `SELECT name FROM slug WHERE id=` + strconv.Itoa(id)

	err := r.DB.QueryRow(
		querySql,
		slug.Name,
	).Scan(&slug.Name)

	if err != nil {
		return slug, err
	}

	return slug, nil
}
