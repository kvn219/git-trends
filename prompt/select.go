package prompt

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"

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
	// fmt.Printf("Filtering repos with %s as the primary programming language.\n", lang)
	return lang
}

// DateTime ..
type DateTime struct {
	Year  string
	Month string
	Day   string
}

// SelectTimeFrame is the time frame for the search.
func SelectTimeFrame() string {
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
	idx, _, err := datePrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed when selecting a programming language", err)
	}
	days := index2Days(idx)
	now := time.Now()
	hours := day2Hour(days)
	then := now.Add(-time.Duration(hours) * time.Hour)
	dt := githubDateFormat(then)
	return dt
}

func githubDateFormat(d time.Time) string {
	dt := DateTime{
		Year:  fourDigitStr(d.Year()),
		Month: twoDigitStr(int(d.Month())),
		Day:   twoDigitStr(d.Day()),
	}
	var tplBuf bytes.Buffer
	const dtTpl = "{{.Year}}-{{.Month}}-{{.Day}}"
	tpl := template.Must(template.New("test").Parse(dtTpl))
	err := tpl.Execute(&tplBuf, dt)
	if err != nil {
		log.Fatal(err)
	}
	return " created:>" + tplBuf.String()
}

func day2Hour(days int) int {
	return days * 24
}

func twoDigitStr(m int) string {
	s := strconv.Itoa(m)
	if len(s) < 2 {
		return fmt.Sprintf("%02d", m)
	}
	return s
}

func fourDigitStr(y int) string {
	s := strconv.Itoa(y)
	if len(s) < 2 {
		return fmt.Sprintf("%04d", y)
	}
	return s
}

func index2Days(idx int) int {
	var days int
	switch idx {
	case 0:
		days = 1800
	case 1:
		days = 14
	case 2:
		days = 31
	case 3:
		days = 180
	case 4:
		days = 360
	default:
		days = 7
	}
	return days
}
