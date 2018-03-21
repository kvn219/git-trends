package ght

import (
	"context"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/kvn219/git-trends/models"
	"github.com/kvn219/git-trends/prompt"
)

// GenerateQueryParams calls several prompts asking the user for search preferences and generates a
// uri params suitable for the github api.
func GenerateQueryParams() string {
	q := prompt.GetKeywords()
	lang := prompt.SelectProgLang()
	dt := prompt.SelectTimeFrame()
	finalQuery := q + " language:" + lang + dt
	return finalQuery
}

// RequestRepos makes a http request to the github api.
func RequestRepos(q string) (*github.RepositoriesSearchResult, *github.Response) {
	ctx := context.Background()
	timeout := time.Duration(5 * time.Second)
	client := github.NewClient(&http.Client{Timeout: timeout})
	opts := &github.SearchOptions{
		Sort:        "stars",
		Order:       "desc",
		ListOptions: github.ListOptions{Page: 0, PerPage: 100},
	}
	output, resp, err := client.Search.Repositories(ctx, q, opts)
	if err != nil {
		return nil, nil
	}
	defer resp.Body.Close()
	return output, resp
}

// ParseRepositories .
func ParseRepositories(output *github.RepositoriesSearchResult) models.Results {
	results := models.Results{}
	for _, r := range output.Repositories {
		rec := models.Record{}
		usr := models.GHUser{}
		rec.ID = r.ID
		rec.Name = r.Name
		rec.URL = r.HTMLURL
		rec.CloneURL = r.CloneURL
		rec.Description = r.Description
		rec.Stars = r.StargazersCount
		rec.ForksCount = r.ForksCount
		rec.CreatedAt = *r.CreatedAt
		usr.Owner = *r.Owner
		rec.Owner = usr
		// fmt.Println(usr.Owner.Login)
		results.Outputs = append(results.Outputs, rec)
	}
	return results
}
