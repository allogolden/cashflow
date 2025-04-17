package postgres

import (
	"database/sql"
)

type AccountModel struct {
	DB *sql.DB
}

func (m AccountModel) CreateAccount(name string, balance float32) (int, error) {
	stmt := `INSERT INTO accounts (name, balance) VALUES (?, ?)`

	result, err := m.DB.Exec(stmt, name, balance)

	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}
