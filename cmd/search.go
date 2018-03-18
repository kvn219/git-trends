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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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
	q = q + " language:go"
	output, resp, err := searchForRepos(q)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	fmt.Printf("Found %d repos!", *output.Total)
	records := parseRepoRecords(output)
	searchResults(records)
	// Write results to json file.
	fpath := getFilePath()
	// serialize record of repos
	serializedRecs := unmarshalRecords(records.Outputs)

	// write results to a json file
	saveRepoResults(fpath, serializedRecs)
	fmt.Println("Finished!!!!")
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
	// for _, repo := range results.Outputs {
	// 	fmt.Println(*repo.ID, *repo.Name, *repo.URL)
	// }
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
		Label:    "What are you searching for?",
		Validate: searchQuery,
		Default:  "",
	}
	q, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user's search query", err)
	}

	langPrompt := promptui.Select{
		Label: "Select Programming Language",
		Items: []string{
			"python",
			"go",
			"javascript",
			"java",
		},
	}

	_, lang, err := langPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when selecting a programming language", err)
	}
	return q + " language:" + lang
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
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Stars | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Stars | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
		--------- Pepper ----------
		{{ "Name:" | faint }}	{{ .Name }}
		{{ "Starts:" | faint }}	{{ .Stars }}
		{{ "Forks:" | faint }}	{{ .ForksCount }}
		`,
	}

	searcher := func(input string, index int) bool {
		repo := dresults[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	listPrompt := promptui.Select{
		Label:     "Stars",
		Items:     dresults,
		Templates: templates,
		Size:      4,
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
	fmt.Println(resp.Request.URL)
	if err != nil {
		return nil, nil, err
	}

	return output, resp, err
}

func deferencePointers(res trends.Results) []trends.UIRecord {
	var out []trends.UIRecord
	for _, repo := range res.Outputs {
		row := trends.UIRecord{
			Name:       *repo.Name,
			URL:        *repo.URL,
			ForksCount: *repo.ForksCount,
			Stars:      *repo.Stars,
		}
		out = append(out, row)
	}
	return out
}
