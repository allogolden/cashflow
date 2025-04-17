package postgres

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) CreateUser(login string, password string, name string, surname string) (int, error) {
	stmt := `INSERT INTO users (login, password, name, surname) VALUES (?, ?, ?, ?)`

	result, err := m.DB.Exec(stmt, login, password, name, surname)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
