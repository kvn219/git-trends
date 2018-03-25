package prompt

import (
	"log"

	"github.com/manifoldco/promptui"
)

// Keywords is prompt that grabs the user's keyword search. For example, a user can search for
// "data science", web, "machine learning", HTTP, etc... Double quotes are necessary if the
// keyword phrase is longer than one word.
func Keywords() string {
	keywordPrompt := promptui.Prompt{
		Label:     "What are you searching for?",
		Validate:  KeywordRequriments,
		AllowEdit: true,
		Default:   "",
	}
	q, err := keywordPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user's search query", err)
	}
	return q
}

// FilePath is prompt that grabs the user's desired output path for all the results of the repo
// search. Currently, the fetch function only limits results to less than 100 repos.
func FilePath() string {
	outPrompt := promptui.Prompt{
		Label:    "Where would you like to save the results?",
		Validate: OutPathRequriments,
		Default:  "",
	}
	fpath, err := outPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when collecting user output path!", err)
	}
	return fpath
}
