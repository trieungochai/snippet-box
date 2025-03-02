package main

import "net/http"

// Because we want this middleware to act on every request that is received,
// we need it to be executed before a request hits our servemux.
// We want the flow of control through our application to look like:
// commonHeaders → servemux → application handler

// Update the signature for the routes() method
// so that it returns a http.Handler instead of *http.ServeMux.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project directory root.
	fileServer := http.FileServer(http.Dir("./ui/static/"))

	// Use the mux.Handle() function to register the file server as the handler for all URL paths that start with "/static/".
	// For matching paths, we strip the "/static" prefix before the request reaches the file server.
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// Swap the route declarations to use the application struct's methods as the handler functions.
	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	// Pass the servemux as the 'next' parameter to the commonHeaders middleware.
	// Because commonHeaders is just a function, and the function returns a http.Handler
	// we don't need to do anything else.
	return commonHeaders(mux)
}
