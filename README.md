# gitsum

[![Build Status](https://cloud.drone.io/api/badges/ctII/gitsum/status.svg)](https://cloud.drone.io/ctII/gitsum)
[![Go Reference](https://pkg.go.dev/badge/github.com/ctII/gitsum.svg)](https://pkg.go.dev/github.com/ctII/gitsum)
[![Go Report Card](https://goreportcard.com/badge/github.com/ctII/gitsum)](https://goreportcard.com/report/github.com/ctII/gitsum)

gitsum is a cli tool for finding information about large number of git repositories in one shot, fairly quickly.

largely built for the personal need to save work before switching computers

## Installing

### Requirements
```git```


```
go install github.com/ctii/gitsum@latest
```

## Usage
```
Usage: gitsum <flags> [path]
  -help
        show help message
  -v    be verbose
```

## Notes
Currently gitsum:
* checks every folder in the current directory recursively, skipping subdirectories of a valid git repo.
* reports a repo as modified if: contains untracked files or changes in current branch || a branch has unpushed commits (this triggers on branches that have no upstream)
