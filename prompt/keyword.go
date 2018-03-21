package prompt

import (
	"log"

	"github.com/kvn219/git-trends/prompt/helpers"
	"github.com/manifoldco/promptui"
)

// GetKeywords is prompt that grabs the user's keyword search. For example, a user can search for
// "data science", web, "machine learning", http, etc... Double quotes are necessary if the
// keyword phrase is longer than one word.
func GetKeywords() string {
	prompt := promptui.Prompt{
		Label:     "What are you searching for?",
		Validate:  helpers.KeywordRequriments,
		AllowEdit: true,
		Default:   "",
	}
	q, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user's search query", err)
	}
	return q
}

// GetFilePath is prompt that grabs the user's desired output path for all the results of the repo
// search. Currently the fetch function only limits results to less than 100 repos.
func GetFilePath() string {
	prompt := promptui.Prompt{
		Label:    "Where would you like to save the results?",
		Validate: helpers.OutPathRequriments,
		Default:  "",
	}
	fpath, err := prompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user output path!", err)
	}
	return fpath
}
