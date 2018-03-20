package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"strconv"
	"time"
)

// DateFormat ..
type DateFormat struct {
	Year  string
	Month string
	Day   string
}

// GithubDateFormat .
func GithubDateFormat(dt time.Time) string {
	df := DateFormat{
		Year:  fourDigitStr(dt.Year()),
		Month: twoDigitStr(int(dt.Month())),
		Day:   twoDigitStr(dt.Day()),
	}
	var tplBuf bytes.Buffer
	const dtTpl = "{{.Year}}-{{.Month}}-{{.Day}}"
	tpl := template.Must(template.New("test").Parse(dtTpl))
	err := tpl.Execute(&tplBuf, df)
	if err != nil {
		log.Fatal(err)
	}
	return " created:>" + tplBuf.String()
}

// Day2Hour .
func Day2Hour(days int) int {
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

// Index2Days .
func Index2Days(idx int) int {
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
