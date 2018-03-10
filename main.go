package main

import (
	"fmt"
	"os"

	"github.com/kvn219/git-trends/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
