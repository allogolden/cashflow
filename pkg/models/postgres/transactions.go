package postgres

import (
	"database/sql"
	"errors"
	"time"

	"golangs.org/snippetbox/pkg/models"
)

type TransactionModel struct {
	DB *sql.DB
}

func (m TransactionModel) CreateTransaction(amount float32, time time.Time, user UserModel, category CategoryModel, account AccountModel) (int, error) {
	stmt := `INSERT INTO transactions (amount, time, user, category, account) VALUES (?, UTC_TIMESTAMP(), ?, ?, ?)`

	result, err := m.DB.Exec(stmt, amount, time, user, category, account)
	
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (m *TransactionModel) GetTransaction(id int) (*models.Transaction, error) {
	stmt := `SELECT id, amount, user, time, category, account FROM transactions WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)

	s := &models.Transaction{}

	err := row.Scan(&s.ID, &s.Amount, &s.Time, &s.User, &s.Category, &s.Account)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *TransactionModel) Latest(user UserModel) ([]*models.Transaction, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE user = ? ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var transactions []*models.Transaction

	for rows.Next() {
		s := &models.Transaction{}
		err = rows.Scan(&s.ID, &s.Amount, &s.Time, &s.User, &s.Category, &s.Account)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
