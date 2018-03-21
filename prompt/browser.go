package prompt

import (
	"fmt"
	"log"
	"strings"

	"github.com/skratchdot/open-golang/open"

	"github.com/kvn219/git-trends/models"
	"github.com/manifoldco/promptui"
)

// BrowserResults .
func BrowserResults(results models.Results) {
	dresults := deferencePointers(results)
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336 {{ .Name | cyan }} ({{ .Stars | red }})",
		Inactive: "  {{ .Name | cyan }} ({{ .Stars | red }})",
		Selected: "\U0001F336 {{ .Name | red | cyan }}",
		Details: `
--------- Repository Information ----------
{{ "Name:" | faint }}	{{ .Name }}
{{ "Stars:" | faint }}	{{ .Stars }}
{{ "Forks:" | faint }}	{{ .ForksCount }}
{{ "Created:" | faint }}	{{ .Created }}
{{ "Description:" | faint }}	{{ .Description }}





















`,
	}

	searcher := func(input string, index int) bool {
		repo := dresults[index]
		name := strings.Replace(strings.ToLower(repo.Name), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)
		return strings.Contains(name, input)
	}

	listPrompt := promptui.Select{
		Label:     "----------List of Results----------",
		Items:     dresults,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}
	i, _, err := listPrompt.Run()
	if err != nil {
		log.Fatal("Prompt failed", err)
	}
	fmt.Printf("You choose number %d: %s\n", i+1, dresults[i].Name)
	open.Run(dresults[i].URL)
}

func deferencePointers(res models.Results) []models.UIRecord {
	var out []models.UIRecord
	for _, repo := range res.Outputs {
		row := models.UIRecord{
			Name:        *repo.Name,
			URL:         *repo.URL,
			ForksCount:  *repo.ForksCount,
			Stars:       *repo.Stars,
			Created:     repo.CreatedAt,
			Description: repo.Description,
		}
		out = append(out, row)
	}
	return out
}
