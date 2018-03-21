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
	"fmt"

	"github.com/kvn219/git-trends/ght"
	"github.com/kvn219/git-trends/prompt"
	"github.com/spf13/cobra"
)

// browseCmd represents the browse command
var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "A brief description of your command",
	Long:  ``,
	Args:  cobra.NoArgs,
	Run:   addBrowse,
}

func init() {
	rootCmd.AddCommand(browseCmd)
}

func addBrowse(*cobra.Command, []string) {
	uri := ght.GenerateQuery()
	output, resp := ght.RequestRepos(uri)
	fmt.Println(resp.Request.URL)
	results := ght.ParseRepositories(output)
	prompt.BrowserResults(results)
}
