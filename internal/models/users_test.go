package models

import (
	"snippetbox.alexedwards.net/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if -short flag is passed
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Define test cases
	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: 2,
			want:   false,
		},
	}

	// Loop through the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use test database (automatically sets up and tears down)
			db := newTestDB(t)

			// Create UserModel with the test DB
			m := UserModel{db}

			// Call the Exists() method
			exists, err := m.Exists(tt.userID)

			// Check results
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
