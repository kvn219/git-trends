package prompt

import (
	"fmt"
	"log"
	"time"

	"github.com/kvn219/git-trends/helpers"
	"github.com/manifoldco/promptui"
)

// SelectProgLang .
func SelectProgLang() string {
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
	return lang
}

// SelectDate .
func SelectDate() string {
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
	days := helpers.Index2Days(idx)
	now := time.Now()
	hours := helpers.Day2Hour(days)
	then := now.Add(-time.Duration(hours) * time.Hour)
	dt := helpers.GithubDateFormat(then)
	return dt
}
