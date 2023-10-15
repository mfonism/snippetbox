package main

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mfonism/snippetbox/internal/models"

	"github.com/julienschmidt/httprouter"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = snippetCreateForm{}

	app.render(w, http.StatusOK, "create.tmpl", data)
}

type snippetCreateForm struct {
	Title string
	Content string
	Expires int
	FieldErrors map[string]string
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := snippetCreateForm{
		FieldErrors: map[string]string{},
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		form.FieldErrors["expires"] = "This field cannot be blank"
	} else {
		expirationChoices := []int{1, 7, 365}
		if !slices.Contains(expirationChoices, expires) {
			form.FieldErrors["expires"] = "This field must equal 1, 7 or 365"
		}
	}
	form.Expires = expires

	title := r.PostForm.Get("title")
	if strings.TrimSpace(title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}
	form.Title = title

	content := r.PostForm.Get("content")
	if strings.TrimSpace(content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}
	form.Content = content

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	id, err := app.snippets.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
