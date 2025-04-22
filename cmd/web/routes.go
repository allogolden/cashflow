package main

import (
	"net/http"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet", app.showSnippet)
	mux.HandleFunc("/snippet/create", app.createSnippet)
	mux.HandleFunc("/api/v1/users/create", app.createUser)
	mux.HandleFunc("/api/v1/login", app.login)
	mux.HandleFunc("/api/v1/accounts/create", app.createAccount)
	mux.HandleFunc("/api/v1/accounts", app.getAccount)
	mux.HandleFunc("/api/v1/user/accounts", app.getUserAccounts)
	mux.HandleFunc("/api/v1/categories/create", app.createCategory)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	return mux
}
