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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/kvn219/git-trends/ght"
	"github.com/kvn219/git-trends/models"
	"github.com/kvn219/git-trends/prompt"
	"github.com/spf13/cobra"
)

// fetchCmd initializes the search program.
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Args:  cobra.NoArgs,
	Short: "Fetch a list of popular github repos and save it to your local computer.",
	Long: `
	Fetch some trending github repositories!

	Supply a keyword or phrase, preferred programming language, and time frame
	of the repository creation date. We'll generate a URI query with your preferences and
	send it to the github api. Afterwards, the results will be parsed into a json file.
	Provide a path to the desired output location and you're all set!
	`,
	Run: addFetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}

func addFetch(cmd *cobra.Command, args []string) {
	extractTransformLoad()
}

func extractTransformLoad() {
	params := ght.GenerateQueryParams()
	output, _ := ght.RequestRepos(params)
	fmt.Println(params)
	// If output from the request is 0, start over...
	if *output.Total == 0 {
		fmt.Println("I couldn't find any thing. Try again...")
		params = ght.GenerateQueryParams()
		output, _ = ght.RequestRepos(params)
	}
	results := ght.ParseRepositories(output)
	serializedResults := unmarshalResults(results.Outputs)
	fp := prompt.GetFilePath()
	saveRepoResults(fp, serializedResults)
}

func unmarshalResults(records []models.Record) []byte {
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
