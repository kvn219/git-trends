package ght

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kvn219/git-trends/ght/prompt"
	"github.com/manifoldco/promptui"
	"github.com/skratchdot/open-golang/open"
)

// GenerateQueryParams calls several prompts asking the user for search preferences and generates a
// uri params suitable for the github api.
func GenerateQueryParams() string {
	q := prompt.Keywords()
	lang := prompt.ProgLang()
	dt := prompt.TimeFrame()
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

// GHUser is a Github user.
type GHUser struct {
	Owner github.User `json:"user"`
}

// Record .
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

// ParseRepositories .
func ParseRepositories(output *github.RepositoriesSearchResult) Results {
	results := Results{}
	for _, r := range output.Repositories {
		rec := Record{}
		usr := GHUser{}
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

// UIRecord is a type for the UI.
type UIRecord struct {
	Name        string
	Stars       int
	ForksCount  int
	URL         string
	Description *string
	Created     github.Timestamp
}

func deferencePointers(res Results) []UIRecord {
	var out []UIRecord
	for _, repo := range res.Outputs {
		row := UIRecord{
			Name:        *repo.Name,
			URL:         *repo.URL,
			ForksCount:  *repo.ForksCount,
			Stars:       *repo.Stars,
			Created:     repo.CreatedAt,
			Description: repo.Description,
		}
		out = append(out, row)
	}
	return out
}

// BrowseRepos generates and list records.
func BrowseRepos(result Results) {
	res := deferencePointers(result)
	showRecords(res)
}

func showRecords(results []UIRecord) {
	customTemp := generateTemplate()
	searcher := func(input string, index int) bool {
		repo := results[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}
	listPrompt := promptui.Select{
		Label:     "----------List of Results----------",
		Items:     results,
		Templates: customTemp,
		Size:      10,
		Searcher:  searcher,
	}
	i, _, err := listPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed", err)
	}
	open.Run(results[i].URL)
}

func generateTemplate() *promptui.SelectTemplates {
	customTemp := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Stars | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Stars | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Repository Information ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Stars:" | faint }}	{{ .Stars }}
{{ "Forks:" | faint }}	{{ .ForksCount }}
{{ "Created:" | faint }}	{{ .Created }}
{{ "Description:" | faint }}	{{ .Description }}
`,
	}
	return customTemp
}
