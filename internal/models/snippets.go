package models

import (
	"database/sql"
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
func (sm *SnippetModel) Insert(title string, content string, expires_at string) (int, error) {
	queryStmt := `INSERT INTO snippets (title, content, created, expires)
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

// This will return a specific snippet based on its id.
// The current implementation is a placeholder and doesn't actually interact with the db.
// It always returns an empty Snippet and nil.
func (sm *SnippetModel) Get(id string) (Snippet, error) {
	return Snippet{}, nil
}

// This will return the 10 most recently created snippets.
// The current implementation is a placeholder and doesn't actually interact with the db.
// It always returns nil and nil.
func (sm *SnippetModel) Latest() ([]Snippet, error) {
	return nil, nil
}
