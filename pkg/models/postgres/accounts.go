package postgres

import (
	"database/sql"
	"errors"

	"golangs.org/snippetbox/pkg/models"
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

func (m *AccountModel) GetAccount(id int) (*models.Account, error) {
	stmt := `SELECT id, name, balance FROM accounts WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Account{}

	err := row.Scan(&s.ID, &s.Name, &s.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}
