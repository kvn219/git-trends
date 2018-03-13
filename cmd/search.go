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
	"net/http"
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
	ctx := context.Background()
	to := &http.Client{Timeout: 10 * time.Second}
	client := github.NewClient(to)
	// Get search query from the command line.
	searchQuery := func(input string) error {
		if len(input) < 2 {
			return errors.New("Search query must have more than 2 characters")
		}
		return nil
	}
	// Set up prompt for the user.
	prompt := promptui.Prompt{
		Label:    "What are you searching for?",
		Validate: searchQuery,
		Default:  "",
	}
	// Grab user's query.
	q, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	// Search repos and sort by stars.
	opts := &github.SearchOptions{Sort: "stars", Order: "desc"}
	output, resp, err := client.Search.Repositories(ctx, q, opts)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	//  Parse results from the body.
	var results = trends.Results{}
	results = parseRepoResults(output, results)
	for _, repo := range results.Outputs {
		// fmt.Printf("%#v\n", repo)
		fmt.Println(repo.ID, repo.Name)
	}
	// Write results to json file.
	fpath, err := getFilePath()
	err = saveRepoResults(fpath, results.Outputs)
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Println("Done!")
}

func parseRepoResults(output *github.RepositoriesSearchResult, results trends.Results) trends.Results {
	for _, r := range output.Repositories {
		rec := trends.Record{}
		rec.ID = *r.ID
		rec.Name = *r.Name
		rec.URL = *r.HTMLURL
		rec.CloneURL = *r.CloneURL
		rec.Description = *r.Description
		rec.Stars = *r.StargazersCount
		results.Outputs = append(results.Outputs, rec)
	}
	return results
}

func saveRepoResults(filename string, results []trends.Record) error {
	b, err := json.Marshal(results)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getFilePath() (string, error) {
	// Get output file path.
	userInput := func(input string) error {
		if len(input) < 2 {
			return errors.New("Query must be longer than 2 characters")
		}
		return nil
	}
	// Set up prompt for the user.
	outputPrompt := promptui.Prompt{
		Label:     "Where would you like to save the results?",
		AllowEdit: true,
		Validate:  userInput,
		Default:   "output.json",
	}
	// Grab user's query.
	fpath, err := outputPrompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return "", err
	}
	return fpath, nil
}
