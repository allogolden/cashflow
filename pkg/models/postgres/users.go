package postgres

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

var ErrInvalidCredentials = errors.New("invalid login or password")
var (
	ErrDuplicateLogin = errors.New("user with this login already exists")
)
var ErrBadDBRequest = errors.New("DB fucked up again")

type UserModel struct {
	DB *sql.DB
}

func (m UserModel) CreateUser(login string, password string, name string, surname string) (int, error) {
	stmt := `INSERT INTO users (login, password, name, surname) VALUES ($1, $2, $3, $4) RETURNING id`

	var id int
	// запускаем QueryRow и сканируем возвращённый id
	err := m.DB.QueryRow(stmt, login, password, name, surname).Scan(&id)
	if err != nil {
		// ловим pq‑ошибку
		if pgErr, ok := err.(*pq.Error); ok {
			// 23505 — duplicate key violation
			if pgErr.Code == "23505" && pgErr.Constraint == "users_login_key" {
				return 0, ErrDuplicateLogin
			}
		}
		return 0, err
	}

	return id, nil
}

func (m UserModel) Login(login string, password string) (error, int) {
	stmt := `SELECT password, id FROM users WHERE login = $1`
	var passDB string
	var idDB int
	err := m.DB.QueryRow(stmt, login).Scan(&passDB, &idDB)

	if err != nil {
		return ErrBadDBRequest, 0
	}
	if passDB != password {
		return ErrInvalidCredentials, 0
	}

	return nil, idDB
}
