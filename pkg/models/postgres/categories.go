package postgres

import (
	"database/sql"
)

type CategoryModel struct {
	DB *sql.DB
}

func (m CategoryModel) CreateCateory(name string) (int, error) {
	stmt := `INSERT INTO users (name) VALUES (?)`

	result, err := m.DB.Exec(stmt, name)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
