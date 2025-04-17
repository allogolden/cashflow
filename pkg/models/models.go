package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: object not found")

type Snippet struct {
	ID int
	Title string
	Content string
	Created time.Time
	Expires time.Time
}

type Transaction struct {
	ID int
	Amount float32
	Time time.Time
	User User
	Category Category
	Account Account
}

type User struct {
	ID int
	name string
	surname string
	login string
	password string
}

type Category struct {
	ID int
	name string
}

type Account struct {
	ID int
	name string
	balance float32
}

type Plan struct {
	ID int
	name string
	amount float32
}
