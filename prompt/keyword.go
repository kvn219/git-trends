package prompt

import (
	"fmt"
	"log"

	"github.com/kvn219/git-trends/helpers"
	"github.com/manifoldco/promptui"
)

// GetKeywords .
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
	fmt.Printf("Filtering repos with %s as the key word.\n", q)
	return q
}

// GetFilePath .
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
