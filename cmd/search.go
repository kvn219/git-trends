// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
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
	"os"

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
 	$ git-trends --search python
	$ git-trends --search golang
	`,
	Run: addSearch,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func addSearch(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	client := github.NewClient(nil)
	var results = trends.Results{}
	validate := func(input string) error {
		if len(input) < 3 {
			return errors.New("Username must have more than 3 characters")
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    "What are you searching for?",
		Validate: validate,
		Default:  "python",
	}

	q, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}
	opts := &github.SearchOptions{Sort: "stars", Order: "desc"}
	output, resp, err := client.Search.Repositories(ctx, q, opts)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	results = parseRepoResults(output, results)
	// for _, q := range args {
	// 	opts := &github.SearchOptions{Sort: "stars", Order: "desc"}
	// 	output, resp, err := client.Search.Repositories(ctx, q, opts)
	// 	defer resp.Body.Close()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(2)
	// 	}
	// 	results = parseRepoResults(output, results)
	// }
	for _, repo := range results.Outputs {
		// fmt.Printf("%#v\n", repo)
		fmt.Println(repo.ID, repo.Name)
		fmt.Println(aGlobalFlag)
	}
	fpath := "output.json"
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
