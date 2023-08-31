package postgresql

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

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

func (r *PostgresqlRep) AddSlugsUser(id int, slugs []models.Slug, duration time.Duration) error {
	var endDate time.Time
	var err error
	if duration != 0 {
		endDate = time.Now().Add(duration)
	}

	for _, slug := range slugs {
		if duration != 0 {
			querySql := `INSERT INTO user_slug (user_id, slug_name, dt_end) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING;`

			_, err = r.DB.Exec(
				querySql,
				id,
				slug.Name,
				endDate,
			)

		} else {
			querySql := `INSERT INTO user_slug (user_id, slug_name) VALUES ($1, $2) ON CONFLICT DO NOTHING;`

			_, err = r.DB.Exec(
				querySql,
				id,
				slug.Name,
			)
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresqlRep) GetCntRandomUsers(cntRandomUsers int) ([]int, error) {
	var users []int

	query := fmt.Sprintf("SELECT id FROM avito_user ORDER BY RANDOM() LIMIT %d ;", cntRandomUsers)

	rows, err := r.DB.Query(query)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return users, err
		}
		users = append(users, id)
	}

	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (r *PostgresqlRep) GetCntUsers() (int, error) {
	var count int
	query := "SELECT COUNT(*) FROM avito_user;"
	err := r.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *PostgresqlRep) DeleteSlugsUser(id int, slugs []models.Slug) error {
	for _, slug := range slugs {
		querySql := `DELETE FROM user_slug WHERE user_id = $1 AND slug_name = $2;`

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

func (r *PostgresqlRep) GetUserArchive(id int, date time.Time) ([]models.Archive, error) {
	var archives []models.Archive

	query := "SELECT user_id, slug_name, assignment, dt FROM archive WHERE user_id = $1 AND dt > $2;"

	rows, err := r.DB.Query(
		query,
		id,
		date,
	)
	if err != nil {
		return archives, err
	}
	defer rows.Close()

	for rows.Next() {
		var archive models.Archive
		if err := rows.Scan(
			&archive.UserId,
			&archive.SlugName,
			&archive.Assigment,
			&archive.DT,
		); err != nil {
			return archives, err
		}
		archives = append(archives, archive)
	}

	if err := rows.Err(); err != nil {
		return archives, err
	}
	return archives, nil
}

func (r *PostgresqlRep) AddArchive(id int, slugs []models.Slug, assigment bool) error {
	for _, slug := range slugs {
		querySql := `INSERT INTO archive (user_id, slug_name, assignment, dt) VALUES ($1, $2, $3, $4);`

		_, err := r.DB.Exec(
			querySql,
			id,
			slug.Name,
			assigment,
			time.Now(),
		)

		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresqlRep) GetUserBySlug(slug models.Slug) ([]models.UserSlug, error) {
	var userSlugs []models.UserSlug

	query := "SELECT user_id, slug_name, dt_end FROM user_slug WHERE slug_name = $1;"

	rows, err := r.DB.Query(
		query,
		slug.Name,
	)
	if err != nil {
		return userSlugs, err
	}
	defer rows.Close()

	for rows.Next() {
		var userSlug models.UserSlug
		if err := rows.Scan(
			&userSlug.UserId,
			&userSlug.SlugId,
			&userSlug.DTEnd,
		); err != nil {
			return userSlugs, err
		}
		userSlugs = append(userSlugs, userSlug)
	}

	if err := rows.Err(); err != nil {
		return userSlugs, err
	}
	return userSlugs, nil
}
