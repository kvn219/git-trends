package ght

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestRepos(t *testing.T) {
	params := "data language:python created:>2017-09-22"
	_, req := RequestRepos(params)
	if req.Response.StatusCode != 200 {
		fmt.Println("Request status should be 200.")
		t.Fail()
	}
}

func TestParseRepositories(t *testing.T) {
	params := "data language:python created:>2017-09-22"
	searchResults, _ := RequestRepos(params)
	results := ParseRepositories(searchResults)
	if len(results.Outputs) == 0 {
		fmt.Println("Should have more than 0 records.")
		t.Fail()
	}
}

func TestRepo(t *testing.T) {
	name1 := "Susy Queue"
	name2 := "Kevin Nguyen"
	usr := &Record{Name: &name1}
	assert.Equal(t, name1, *usr.Name)
	assert.NotEqual(t, name2, *usr.Name)
}
