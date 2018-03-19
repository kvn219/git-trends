// Copyright © 2018 Kevin Nguyen kvn219@nyu.edu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kvn219/git-trends/trends"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var priority int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Args:  cobra.NoArgs,
	Short: "Your search query!",
	Long: `You can search for any thing! For example:

	$ git-trends search
	 ✗ What are you searching for?: python
	$ git-trends user # display user info
	$ git-trends repo # display repo info
	`,
	Run: addSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func addSearch(cmd *cobra.Command, args []string) {
	q := grabUserQuery()
	output, resp, err := searchForRepos(q)
	if err != nil {
		log.Fatal("Query failed to get response.", err)
	}
	defer resp.Body.Close()
	fmt.Println("Search Query URL:")
	fmt.Println(resp.Request.URL.String())
	records := parseRepoRecords(output)
	searchResults(records)
	// Write results to json file.
	fpath := getFilePath()
	// serialize record of repos
	serializedRecs := unmarshalRecords(records.Outputs)
	// write results to a json file
	saveRepoResults(fpath, serializedRecs)
}

func parseRepoRecords(output *github.RepositoriesSearchResult) trends.Results {
	results := trends.Results{}
	for _, r := range output.Repositories {
		rec := trends.Record{}
		rec.ID = r.ID
		rec.Name = r.Name
		rec.URL = r.HTMLURL
		rec.CloneURL = r.CloneURL
		rec.Description = r.Description
		rec.Stars = r.StargazersCount
		rec.ForksCount = r.ForksCount
		results.Outputs = append(results.Outputs, rec)
	}
	return results
}

func unmarshalRecords(records []trends.Record) []byte {
	b, err := json.Marshal(records)
	if err != nil {
		log.Fatal("Marshalling failed", err)
	}
	return b
}

func saveRepoResults(filename string, results []byte) error {
	err := ioutil.WriteFile(filename, results, 0644)
	if err != nil {
		log.Fatal("Writing file failed", err)
	}
	return nil
}

func getFilePath() string {
	prompt := promptui.Prompt{
		Label:    "Where would you like to save the results?",
		Validate: reqOutputPath,
		Default:  "",
	}
	fpath, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user output path!", err)
	}
	return fpath
}

func reqOutputPath(input string) error {
	if len(input) < 2 {
		return errors.New("Search query must have more than 2 characters")
	}
	if filepath.Ext(input) != ".json" {
		return errors.New("Must be json file")
	}
	return nil
}

func grabUserQuery() string {
	prompt := promptui.Prompt{
		Label:     "What are you searching for?",
		Validate:  searchQuery,
		AllowEdit: true,
		Default:   "",
	}
	q, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user's search query", err)
	}
	fmt.Printf("Filtering repos with %s as the key word.\n", q)
	langPrompt := promptui.Select{
		Label: "Select Programming Language",
		Items: []string{
			"python",
			"go",
			"javascript",
			"java",
			"juptyer notebook",
		},
	}
	_, lang, err := langPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when selecting a programming language", err)
	}
	fmt.Printf("Filtering repos with %s as the primary programming language.\n", lang)
	cd := filterByDate()
	finalQuery := q + " language:" + lang + cd
	fmt.Println(finalQuery)
	return finalQuery
}

func searchQuery(input string) error {
	if len(input) < 2 {
		return errors.New("Search query must have more than 2 characters")
	}
	return nil
}

func searchResults(results trends.Results) {
	dresults := deferencePointers(results)
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Stars | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Stars | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Repository Information ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Stars:" | faint }}	{{ .Stars }}
{{ "Forks:" | faint }}	{{ .ForksCount }}
{{ "Created:" | faint }}	{{ .CreatedAt }}
{{ "Description:" | faint }}	{{ .Description }}











`,
	}

	searcher := func(input string, index int) bool {
		repo := dresults[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	listPrompt := promptui.Select{
		Label:     "----------List of Results----------",
		Items:     dresults,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}
	i, _, err := listPrompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	fmt.Printf("You choose number %d: %s\n", i+1, dresults[i].Name)
}

func searchForRepos(q string) (*github.RepositoriesSearchResult, *github.Response, error) {
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
		return nil, nil, err
	}

	return output, resp, err
}

func deferencePointers(res trends.Results) []trends.UIRecord {
	var out []trends.UIRecord
	for _, repo := range res.Outputs {
		row := trends.UIRecord{
			Name:        *repo.Name,
			URL:         *repo.URL,
			ForksCount:  *repo.ForksCount,
			Stars:       *repo.Stars,
			Description: repo.Description,
			CreatedAt:   repo.CreatedAt,
		}
		out = append(out, row)
	}
	return out
}

func filterByDate() string {
	datePrompt := promptui.Select{
		Label: "How far do you want to go back?",
		Items: []string{
			"All time",
			"Two weeks",
			"A month",
			"Six months",
			"The past year",
		},
	}
	idx, date, err := datePrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when selecting a programming language", err)
	}
	fmt.Println("Filtering repos created within " + date)
	now := time.Now()
	var days int
	if idx == 0 {
		days = 1800
	} else if idx == 1 {
		days = 14
	} else if idx == 2 {
		days = 31
	} else if idx == 3 {
		days = 180
	} else if idx == 4 {
		days = 360
	} else {
		days = 7
	}
	hours := day2Hour(days)
	then := now.Add(-time.Duration(hours) * time.Hour)
	// thenFMT := then.Format(time.RFC1123)
	cd := githubDateFormat(then)
	return cd
}

func day2Hour(days int) int {
	return days * 24
}

// Date ..
type Date struct {
	Year  string
	Month string
	Day   string
}

func githubDateFormat(dt time.Time) string {
	y := enforce4Digit(dt.Year())
	m := enforce2Digit(int(dt.Month()))
	d := enforce2Digit(dt.Day())
	s := Date{
		Year:  y,
		Month: m,
		Day:   d,
	}
	var tpl bytes.Buffer
	const textTmp = "{{.Year}}-{{.Month}}-{{.Day}}"
	tmpl := template.Must(template.New("test").Parse(textTmp))
	err := tmpl.Execute(&tpl, s)
	if err != nil {
		log.Fatal(err)
	}
	return " created:>" + tpl.String()
}

func enforce2Digit(m int) string {
	s := strconv.Itoa(m)
	if len(s) < 2 {
		return fmt.Sprintf("%02d", m)
	}
	return s
}

func enforce4Digit(y int) string {
	s := strconv.Itoa(y)
	if len(s) < 2 {
		return fmt.Sprintf("%04d", y)
	}
	return s
}
