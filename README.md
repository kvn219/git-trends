# Git-trends

**A simple command line application to search repositories trending on GitHub.**

Git-trends lets you look through a list of popular repositories, access them on your browser, and save results to your local computer using the [google/go-github](https://github.com/google/go-github) package. I'm writing this app to learn Golang. So feel free to send me a message or submit an issue if you have any tips, suggestions, or comments.

## Getting started

### Clone the repo and run the binary file

```bash
git clone https://github.com/kvn219/git-trends.git && cd git-trends
./git-trends help
```

## Two basic commands

There are two simple commands. The first is `browse`, which walks you through a series of steps to access a list of the latest repositories ordered by the number of stars on GitHub. Afterward, you can look through a list of repositories in the terminal. When you select a repository from the list your default, your browser will open up with the of the repository you chose.

```bash
>>> ./git-trends browse # browse
>>> ✔ What are you searching for?: # "data science"
>>> ? Select Programming Language: # python
>>> ? How far do you want to go back?:  # Last month...
```

The second command `fetch`, also walks you through a series of steps to access a list of latest repositories ordered by the number of stars on GitHub. After fetching some repos, you can save the output
to a JSON file.

```bash
>>> ./git-trends fetch # fetch
>>> ✔ What are you searching for?: # "dev ops"
>>> ? Select Programming Language: # go
>>> ? How far do you want to go back?:  # Last year...
>>> ✔ Where would you like to save the results?: data-science-golang.json
```

## Features

* [x] Search repos
* [x] Select repo and open browser
* [x] Save output to local machine
* [ ] Add some more tests
* [ ] Save to database
* [ ] Improve UI
* [ ] Add Travis CI
* [ ] Add Docker

## Project Layout

```bash
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   ├── browse.go
│   ├── fetch.go
│   └── root.go
├── ght
│   └── ght.go
├── main.go
├── models
│   └── Repo.go
└── prompt
    ├── browser.go
    ├── keyword.go
    ├── select.go
    └── validators.go
```
