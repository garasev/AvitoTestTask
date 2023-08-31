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
	querySql := `INSERT INTO slug (name) VALUES ($1) RETURNING id;`

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
	querySql := `SELECT name FROM slug WHERE id=` + strconv.Itoa(id) + `;`

	err := r.DB.QueryRow(querySql).Scan(&slug.Name)

	if err != nil {
		return slug, err
	}

	return slug, nil
}

func (r *PostgresqlRep) GetSlugs() ([]models.Slug, error) {
	var slugs []models.Slug
	querySql := `SELECT name FROM slug;`

	rows, err := r.DB.Query(querySql)

	if err != nil {
		return slugs, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug models.Slug
		err := rows.Scan(
			&slug.Name,
		)
		if err != nil {
			return slugs, err
		}
		slugs = append(slugs, slug)
	}

	if err = rows.Err(); err != nil {
		return slugs, err
	}

	return slugs, nil
}

func (r *PostgresqlRep) DeleteSlug(slug models.Slug) error {
	querySql := `DELETE FROM slug WHERE name =$1;`

	_, err := r.DB.Exec(querySql, slug.Name)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresqlRep) GetUsers() ([]models.AvitoUser, error) {
	var users []models.AvitoUser
	querySql := `SELECT id FROM avito_user;`

	rows, err := r.DB.Query(querySql)

	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.AvitoUser
		err := rows.Scan(
			&user.Id,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return users, err
	}

	return users, nil
}

func (r *PostgresqlRep) GetUserSlugs(id int) ([]models.Slug, error) {
	var slugs []models.Slug
	querySql := `SELECT s.name FROM user_slug AS us JOIN slug AS s ON us.slug_id = s.id WHERE us.dt_end > NOW() OR us.dt_end IS NULL;`

	rows, err := r.DB.Query(querySql)

	if err != nil {
		return slugs, err
	}
	defer rows.Close()

	for rows.Next() {
		var slug models.Slug
		err := rows.Scan(
			&slug.Name,
		)
		if err != nil {
			return slugs, err
		}
		slugs = append(slugs, slug)
	}

	if err = rows.Err(); err != nil {
		return slugs, err
	}

	return slugs, nil
}

func (r *PostgresqlRep) AddUsers(cnt int) error {
	for i := 0; i < cnt; i++ {
		querySql := `INSERT INTO avito_user DEFAULT VALUES;`

		_, err := r.DB.Exec(querySql)

		if err != nil {
			return err
		}
	}
	return nil
}
