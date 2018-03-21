package models

import "github.com/google/go-github/github"

// GHUser is a Github user.
type GHUser struct {
	Owner github.User `json:"user"`
}
