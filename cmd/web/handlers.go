package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"snippetbox.t10i.net/internal/models"
)

// Define a snippetCreateForm struct to represent the form data and validation errors for the form fields.
// Note that all the struct fields are deliberately exported (i.e. start with a capital letter).
// This is because struct fields must be exported in order to be read by the html/template package
// when rendering the template.
type snippetCreateForm struct {
	Title       string
	Content     string
	ExpiresAt   int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Call the newTemplateData() helper to get a templateData struct containing
	// the 'default' data (which for now is just the current year),
	// and add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	// Pass the data to the render() helper as normal.
	app.render(w, r, http.StatusOK, "home.tmpl", data)
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.tmpl", data)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Initialize a new createSnippetForm instance and pass it to the template.
	// Notice how this is also a great opportunity to set any default or 'initial' values for the form
	// --- here we set the initial value for the snippet expiry to 365 days.
	data.Form = snippetCreateForm{
		ExpiresAt: 365,
	}

	app.render(w, r, http.StatusOK, "create.tmpl", data)
}

func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// First we call r.ParseForm() which adds any data in POST request bodies to the r.PostForm map.
	// If there are any errors, we use our app.ClientError() helper to
	// send a 400 Bad Request response to the user.
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// The r.PostForm.Get() method always returns the form data as a *string*.
	// However, we're expecting our expires value to be a number,
	// and want to represent it in our Go code as an integer.
	// So we need to manually convert the form data to an integer using strconv.Atoi(),
	// and we send a 400 Bad Request response if the conversion fails.
	expires_at, err := strconv.Atoi(r.PostForm.Get("expires_at"))

	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Create an instance of the snippetCreateForm struct containing the values
	// from the form and an empty map for any validation errors.
	form := snippetCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		ExpiresAt:   expires_at,
		FieldErrors: map[string]string{},
	}

	// Update the validation checks so that they operate on the snippetCreateForm instance.
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This field cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "This field cannot be more than 100 characters long"
	}

	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This field cannot be blank"
	}

	// Check the expires value matches one of the permitted values (1, 7 or 365).
	if form.ExpiresAt != 1 && form.ExpiresAt != 7 && form.ExpiresAt != 365 {
		form.FieldErrors["expires_at"] = "This field must equal 1, 7 or 365"
	}

	// If there are any validation errors, then re-display the create.tmpl template,
	// passing in the snippetCreateForm instance as dynamic data in the Form field.
	// Note that we use the HTTP status code 422 Unprocessable Entity
	// when sending the response to indicate that there was a validation error.
	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.tmpl", data)
		return
	}

	// Pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back.
	id, err := app.snippets.Insert(form.Title, form.Content, form.ExpiresAt)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
