package postgresql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

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

func (r *PostgresqlRep) AddSlug(slug models.Slug) error {
	querySql := `INSERT INTO slug (name) VALUES ($1);`

	_, err := r.DB.Exec(
		querySql,
		slug.Name,
	)

	if err != nil {
		return err
	}

	return nil
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
	querySql := fmt.Sprintf("SELECT slug_name FROM user_slug WHERE user_id = %d AND (dt_end > NOW() OR dt_end IS NULL);", id)

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

func (r *PostgresqlRep) CheckUserExist(id int) (bool, error) {
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM avito_user WHERE %s);", "id ="+strconv.Itoa(id))

	err := r.DB.QueryRow(query).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *PostgresqlRep) CheckSlugsExist(slugs []models.Slug) (bool, error) {
	if len(slugs) == 0 {
		return true, nil
	}

	placeholders := make([]string, len(slugs))
	values := make([]interface{}, len(slugs))

	for i, slug := range slugs {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		values[i] = slug.Name
	}

	query := fmt.Sprintf("SELECT name FROM slug WHERE name IN (%s)", strings.Join(placeholders, ", "))

	rows, err := r.DB.Query(query, values...)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	existingSlugs := make(map[string]bool)

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return false, err
		}
		existingSlugs[name] = true
	}

	if err := rows.Err(); err != nil {
		return false, err
	}
	for _, slug := range slugs {
		if !existingSlugs[slug.Name] {
			return false, nil
		}
	}

	return true, nil
}

func (r *PostgresqlRep) AddSlugsUser(id int, slugs []models.Slug) error {
	for _, slug := range slugs {
		querySql := `INSERT INTO user_slug (user_id, slug_name) VALUES ($1, $2);`

		_, err := r.DB.Exec(
			querySql,
			id,
			slug.Name,
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresqlRep) DeleteSlugsUser(id int, slugs []models.Slug) error {
	for _, slug := range slugs {
		querySql := `DELETE FROM user_slug WHERE user_id = $1 AND name = $2;`

		_, err := r.DB.Exec(
			querySql,
			id,
			slug.Name,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
