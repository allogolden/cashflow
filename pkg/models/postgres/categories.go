package postgres

import (
	"database/sql"
	"errors"

	"golangs.org/snippetbox/pkg/models"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) CreateCateory(name string, user int64) (int, error) {
	stmt := `INSERT INTO categories (name, user_id) VALUES ($1, $2) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, name, user).Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *CategoryModel) GetCategory(id int) (*models.Category, error) {
	stmt := `SELECT id, name FROM categories WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Category{}

	err := row.Scan(&s.ID, &s.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *CategoryModel) GetUserCategories(user int) ([]*models.Category, error) {
	stmt := `SELECT id, name, balance FROM categories
	WHERE user = ? ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, user)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*models.Category

	for rows.Next() {
		s := &models.Category{}
		err = rows.Scan(&s.ID, &s.Name)
		if err != nil {
			return nil, err
		}

		categories = append(categories, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
