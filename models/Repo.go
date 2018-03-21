package models

import (
	"github.com/google/go-github/github"
)

// Record from git search
type Record struct {
	ID          *int64           `json:"id,omitempty"`
	Name        *string          `json:"name"`
	URL         *string          `json:"url"`
	Description *string          `json:"description"`
	CloneURL    *string          `json:"clone_url"`
	Stars       *int             `json:"stars"`
	ForksCount  *int             `json:"forks_count"`
	CreatedAt   github.Timestamp `json:"created_at"`
	Owner       GHUser           `json:"owner"`
}

// Results from git search
type Results struct {
	Outputs []Record
}

// GHUser .
type GHUser struct {
	Owner github.User `json:"user"`
}

// UIRecord is a type for the UI.
type UIRecord struct {
	Name        string
	Stars       int
	ForksCount  int
	URL         string
	Description *string
	Created     github.Timestamp
}
