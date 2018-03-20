package helpers

import (
	"errors"
	"path/filepath"
)

// KeywordRequriments .
func KeywordRequriments(input string) error {
	if len(input) < 2 {
		return errors.New("Search query must have more than 2 characters")
	}
	return nil
}

// OutPathRequriments ..
func OutPathRequriments(input string) error {
	n := len(input)
	ext := filepath.Ext(input)
	if n < 2 || ext != ".json" {
		return errors.New("Search query must have more than 2 characters and have a .json ext")
	}
	return nil
}
