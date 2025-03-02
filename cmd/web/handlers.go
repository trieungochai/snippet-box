package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"snippetbox.t10i.net/internal/models"
)

// Change the signature of the home handler
// so it is defined as a method against *application.
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "Go")

	snippets, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Initialize a slice containing the paths to the two files.
	// It's important to note that
	// the file containing our base template must be the *first* file in the slice.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/home.tmpl",
	}

	// Use the template.ParseFiles() func to read the files and store the templates in a template set.
	// Notice that we use ... to pass the contents of the files slice as variadic arguments.
	// If there's an error, we log the detailed error message,
	// use the http.Error() function to send an Internal Server Error response to the user,
	// and then return from the handler so no subsequent code is executed.
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		// 	// Because the home handler is now a method against the application struct
		// 	// it can access its fields, including the structured logger.
		// 	// We'll use this to create a log entry at Error level containing the error message,
		// 	// also including the request method and URI as attributes to assist with debugging.
		app.serverError(w, r, err) // Use the server serverError() helper
		return
	}

	// Create an instance of a templateData struct holding the slice of snippets.
	data := templateData{
		Snippets: snippets,
	}

	// Pass in the templateData struct when executing the template.
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// Use the SnippetModel's Get() method to retrieve the data for a specific record based on its ID.
// If no matching record is found, return a 404 Not Found response.
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

	// Initialize a slice containing the paths to the view.tmpl file,
	// plus the base layout and navigation partial that we made earlier.
	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
	}

	// Parse the template files...
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Create an instance of a templateData struct holding the snippet data.
	data := templateData{
		Snippet: snippet,
	}

	// Pass in the templateData struct when executing the template.
	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverError(w, r, err)
	}
}

// Change the signature of the home handler
// so it is defined as a method against *application.
func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a form for creating a new snippet..."))
}

// Change the signature of the home handler
// so it is defined as a method against *application.
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Create some variables holding dummy data. We'll remove these later on during the build.
	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires_at := 7

	// Pass the data to the SnippetModel.Insert() method, receiving the ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires_at)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Redirect the user to the relevant page for the snippet.
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
