package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"snippetbox.t10i.net/internal/models"
)

// Define an application struct to hold the application-wide dependencies for the web app.
// Add a snippets field to the application struct.
// This will allow us to make the SnippetModel object available to our handlers.
// Add a templateCache field to the application struct.
type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	// Define a new command-line flag with the name 'addr', a default value of ":4000"
	// and some short help text explaining what the flag controls.
	// The value of the flag will be stored in the the addr variable at runtime.
	addr := flag.String("addr", ":4000", "HTTP server network address")

	// Define a new command-line flag for the MySQL DSN string.
	dsn := flag.String("dsn", "web:normaluser@/snippetbox?parseTime=true", "MySQL data source name")

	// Importantly, we use the flag.Parse() function to parse the command-line flag.
	// This reads in the command-line flag value and assigns it to the addr variable.
	// You need to call this *before* you use the addr variable
	// otherwise it will always contain the default value of ":4000".
	// If any errors are encountered during parsing the application will be terminated.
	flag.Parse()

	// Use the slog.New() function to initialize a new structured logger,
	// which writes to the standard out stream and uses the default settings.
	loggerHandler := slog.NewTextHandler(os.Stdout, nil)
	logger := slog.New(loggerHandler)

	// To keep the main() function tidy
	// I've put the code for creating a connection pool into the separate openDB() function below.
	// We pass openDB() the DSN from the command-line flag.
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// We also defer a call to db.Close(),
	// so that the connection pool is closed before the main() function exits.
	defer db.Close()

	// Initialize a new template cache.
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Init a new instance of our application struct, containing the dependencies
	// Init a models.SnippetModel instance containing the connection pool and add it to the application dependencies.
	// And add it to the application dependencies.
	app := &application{
		logger:        logger,
		snippets:      &models.SnippetModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	// Call the new app.routes() method to get the servemux containing our routes,
	// and pass that to http.ListenAndServe().
	// Because the err variable is now already declared in the code above, we need
	// to use the assignment operator = here, instead of the := 'declare and assign' operator.
	err = http.ListenAndServe(*addr, app.routes())

	// And we also use the Error() method to log any error message returned by
	// http.ListenAndServe() at Error severity (with no additional attributes),
	// and then call os.Exit(1) to terminate the application with exit code 1.
	logger.Error(err.Error())
	os.Exit(1)
}

// The openDB() function wraps sql.Open()
// and returns a sql.DB connection pool for a given DSN.
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
