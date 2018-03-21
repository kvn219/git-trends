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
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/github"
	"github.com/kvn219/git-trends/models"
	"github.com/kvn219/git-trends/prompt"
	"github.com/spf13/cobra"
)

// searchCmd initializes the search program.
var searchCmd = &cobra.Command{
	Use:   "fetch",
	Args:  cobra.NoArgs,
	Short: "Keyword search query.",
	Long: `Begin the program by entering a keyword search query.

	For example:

	$ git-trends fetch
	$ What are you searching for?: https
	$ git-trends repo # display repo info
	`,
	Run: addSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func addSearch(cmd *cobra.Command, args []string) {
	q := prompt.GetKeywords()
	lang := prompt.SelectProgLang()
	dt := prompt.SelectTimeFrame()
	fq := generateQuery(q, lang, dt)
	output, resp := findRepos(fq)
	fmt.Println("Query: " + resp.Request.URL.String())
	records := parseRepositories(output)
	prompt.BrowserResults(records)
	serializedRecs := unmarshalRecords(records.Outputs)
	fpath := prompt.GetFilePath()
	saveRepoResults(fpath, serializedRecs)
}

func parseRepositories(output *github.RepositoriesSearchResult) models.Results {
	results := models.Results{}
	for _, r := range output.Repositories {
		rec := models.Record{}
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

func unmarshalRecords(records []models.Record) []byte {
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

func generateQuery(q, lang, date string) string {
	finalQuery := q + " language:" + lang + date
	return finalQuery
}

func findRepos(q string) (*github.RepositoriesSearchResult, *github.Response) {
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
