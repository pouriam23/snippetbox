package mocks

import (
	"database/sql"
	"snippetbox.alexedwards.net/internal/models"
	"time"
)

// Dummy snippet to return in tests
var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, sql.ErrNoRows
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
