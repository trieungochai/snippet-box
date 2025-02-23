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
// The current implementation is a placeholder and doesn't actually interact with the db.
// It always returns 0 and nil.
func (sm *SnippetModel) Insert(title string, content string, expiresAt string) (int, error) {
	return 0, nil
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
