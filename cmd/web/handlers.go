package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"golangs.org/snippetbox/pkg/models"
	"golangs.org/snippetbox/pkg/models/postgres"
)

// Обработчик главной странице.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})

}

// Обработчик для отображения содержимого заметки.
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	s, err := app.snippets.Get(id)

	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

}

// Обработчик для создания новой заметки.
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title := "Story of the Snail"
	content := "Lorem"
	expires := "7"

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)

}

func (app *application) createUser(w http.ResponseWriter, r *http.Request) {
	if app.users == nil {
		log.Println("ERROR: app.models.User is nil!")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if app.users.DB == nil {
		log.Println("ERROR: app.models.User.DB is nil!")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
	}
	login := r.FormValue("login")
	password := r.FormValue("password")
	name := r.FormValue("name")
	surname := r.FormValue("surname")

	_, err := app.users.CreateUser(login, password, name, surname)

	if err != nil {
		switch err {
		case postgres.ErrDuplicateLogin:
			http.Error(w, err.Error(), http.StatusConflict)
		default:
			app.serverError(w, err)
		}
		return
	}

}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	if app.users == nil {
		log.Println("ERROR: app.models.User is nil!")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if app.users.DB == nil {
		log.Println("ERROR: app.models.User.DB is nil!")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)

		app.clientError(w, http.StatusMethodNotAllowed)
	}

	login := r.FormValue("login")
	password := r.FormValue("password")
	fmt.Println(login, password)

	err, id := app.users.Login(login, password)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "404",
			"message": "Invalid credentials",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "Login successful",
		"user_id": id, // если вы из Login возвращаете id
	})
}
