package postgres

import (
	"database/sql"
	"errors"

	"golangs.org/snippetbox/pkg/models"
)

type AccountModel struct {
	DB *sql.DB
}

func (m AccountModel) CreateAccount(name string, balance float64, user_id int64) (int, error) {
	stmt := `INSERT INTO accounts (name, balance, user_id) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := m.DB.QueryRow(stmt, name, balance, user_id).Scan(&id)

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *AccountModel) GetAccount(id int64) (*models.Account, error) {
	stmt := `SELECT id, name, balance FROM accounts WHERE id = $1`

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

func (m *AccountModel) GetUserAcccounts(user int64) ([]*models.Account, error) {
	stmt := `SELECT id, name, balance FROM accounts
	WHERE user_id = $1 ORDER BY name DESC LIMIT 10`

	rows, err := m.DB.Query(stmt, user)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var accounts []*models.Account

	for rows.Next() {
		s := &models.Account{}
		err = rows.Scan(&s.ID, &s.Name, &s.Balance)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
