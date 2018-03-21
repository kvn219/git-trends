# Git-trends

**A simple command line application to search repositories trending on github.**

Git-trends is a simple command line application to search repositories trending on github! You can
look through a list of popular repos, access them on your browser, and save results to your local computer.

## Features

* [x] Search (Repos)
* [x] Select Repo and Open Browser
* [x] Save (output directory)
* [ ] Improve UI

```bash
>>>./git-trends fetch
>>>✔ What are you searching for?: # "data science"
>>>? Select Programming Language: # python
>>>? How far do you want to go back?:  # Last month...
```

## Project Layout

```bash
├── LICENSE
├── Makefile
├── README.md
├── cmd
│   ├── fetch.go
│   └── root.go
├── main.go
├── models
│   └── Repo.go
└── prompt
    ├── browser.go
    ├── helpers
    │   └── validators.go
    ├── keyword.go
    └── select.go
```
