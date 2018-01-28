package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/github"
)

type Repo struct {
	HTMLURL string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Simple cli")
	fmt.Print("User: ")
	u, _ := reader.ReadString('\n')
	user := strings.Replace(u, "\n", "", -1)

	fmt.Print("Repo: ")
	r, _ := reader.ReadString('\n')
	repo := strings.Replace(r, "\n", "", -1)

	info := fmt.Sprintf("Finding %s repo from %s.", repo, user)
	fmt.Println(info)

	ctx := context.Background()
	client := github.NewClient(nil)
	output, _, err := client.Repositories.Get(ctx, user, repo)

	if err != nil {
		fmt.Println(err)
	}

	name := *output.HTMLURL
	url := "Go to link: %s"
	message := fmt.Sprintf(url, name)
	fmt.Println(message)
}
