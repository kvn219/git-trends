# Git-trends

**A simple command line application to search repositories trending on github.**

Git-trends lets you look through a list of popular repositories, access them on your browser, and save results to your local computer using the [google/go-github](https://github.com/google/go-github) package.

## Getting started

### Clone the repo and run the binary file

```bash
git clone https://github.com/kvn219/git-trends.git && cd git-trends
./git-trends help
```

## Two basic commands!

There are two simple commands. The first is `browse` which lets you look through
a list of repositories in the command line. When you select a repository from the list your default, your browser will open up with the github link of repository you selected.

```bash
>>> ./git-trends browse
>>> ✔ What are you searching for?: # "data science"
>>> ? Select Programming Language: # python
>>> ? How far do you want to go back?:  # Last month...
```

The second is `fetch`, which walks you through a serise of steps to access a list of latest repositories ordered by number of stars on github.

```bash
>>> ./git-trends fetch
>>> ✔ What are you searching for?: # "dev ops"
>>> ? Select Programming Language: # go
>>> ? How far do you want to go back?:  # Last month...
>>> ✔ Where would you like to save the results?: data-science-golang.json
```

## Features

* [x] Search repos
* [x] Select repo and open browser
* [x] Save output to local machine
* [ ] Save to database
* [ ] Improve UI

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
