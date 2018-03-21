package ght

import (
	"fmt"
	"testing"
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
