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

// SelectProgLang is a prompt for selecting the programming language
// you wish condition the search.
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
	return lang
}

// DateTime is a date format need to pass into the query template.
type DateTime struct {
	Year  string
	Month string
	Day   string
}

// SelectTimeFrame is a prompt for selecting the extent of time you want
// to go back during the search. It's based on created date, for example: "Last year..." means
// we're going to look for repos created within the last year.
func SelectTimeFrame() string {
	datePrompt := promptui.Select{
		Label: "How far do you want to go back? (based on created date)",
		Items: []string{
			"Way back... (5 years)",
			"Last year...",
			"Last 6 months...",
			"Last month...",
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
		days = 360
	case 2:
		days = 180
	case 3:
		days = 31
	default:
		days = 1800
	}
	return days
}
