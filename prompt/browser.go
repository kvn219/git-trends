package prompt

import (
	"fmt"
	"log"
	"strings"

	"github.com/kvn219/git-trends/search"
	"github.com/manifoldco/promptui"
)

// BrowserResults .
func BrowserResults(results search.Results) {
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
{{ "Created:" | faint }}	{{ .CreatedAt }}
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
}

func deferencePointers(res search.Results) []search.UIRecord {
	var out []search.UIRecord
	for _, repo := range res.Outputs {
		row := search.UIRecord{
			Name:        *repo.Name,
			URL:         *repo.URL,
			ForksCount:  *repo.ForksCount,
			Stars:       *repo.Stars,
			Description: repo.Description,
			CreatedAt:   repo.CreatedAt,
		}
		out = append(out, row)
	}
	return out
}
