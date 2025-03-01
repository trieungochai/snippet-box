package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a Snippet type to hold the data for an individual snippet.
// Notice how the fields of the struct correspond to the fields in our MySQL snippets table.
type Snippet struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	ExpiresAt time.Time
}

// Define a SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
// Returns:
// 1. The ID of the newly inserted snippet (integer).
// 2. An error if something goes wrong.
func (sm *SnippetModel) Insert(title string, content string, expires_at int) (int, error) {
	queryStmt := `INSERT INTO snippets (title, content, created_at, expires_at)
    VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// Use the Exec() method on the embedded connection pool to execute the statement.
	// The first parameter is the SQL statement,
	// followed by the values for the placeholder parameters: title, content and expiry in that order.
	// This method returns a sql.Result type, which contains some
	// basic information about what happened when the statement was executed.
	sqlResult, err := sm.DB.Exec(queryStmt, title, content, expires_at)
	if err != nil {
		return 0, err
	}

	// Use the LastInsertId() method on the result to get the ID of our
	// newly inserted record in the snippets table.
	id, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	// The ID returned has the type int64, so we convert it to an int type before returning.
	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (Snippet, error) {
	queryStmt := `SELECT id, title, content, created_at, expires_at FROM snippets
	WHERE expires_at > UTC_TIMESTAMP() and id = ?`
	// Use the QueryRow() method on the connection pool to execute our SQL statement,
	// passing in the untrusted id variable as the value for the placeholder param.
	// This returns a pointer to a sql.Row object which holds the result from the db.
	row := sm.DB.QueryRow(queryStmt, id)

	// Init a new zeroed Snippet struct.
	var snippet Snippet

	// Use row.Scan() to copy the values from each field in sql.Row to the corresponding field in the Snippet struct.
	// Notice that the arguments to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of columns returned by your statement.
	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.CreatedAt, &snippet.ExpiresAt)

	// If the query returns no rows, then row.Scan() will return a sql.ErrNoRows error.
	// We use the errors.Is() function check for that error specifically,
	// and return our own ErrNoRecord error instead.
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	// If everything went OK, then return the filled Snippet struct.
	return snippet, nil
}

// This will return the 10 most recently created snippets.
// The current implementation is a placeholder and doesn't actually interact with the db.
// It always returns nil and nil.
func (sm *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
