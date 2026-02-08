package model

import (
	"time"
)

// Entry represents a single password entry in the vault.
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	URL       string    `json:"url,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Vault represents the decrypted content of the password store.
type Vault struct {
	Entries []Entry `json:"entries"`
}
