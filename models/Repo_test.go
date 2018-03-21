package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	name1 := "Susy Queue"
	name2 := "Kevin Nguyen"
	usr := &Record{Name: &name1}
	assert.Equal(t, name1, *usr.Name)
	assert.NotEqual(t, name2, *usr.Name)
}
