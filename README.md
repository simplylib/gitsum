# gitsum

[![Build Status](https://cloud.drone.io/api/badges/simplylib/gitsum/status.svg)](https://cloud.drone.io/simplylib/gitsum)
[![Go Reference](https://pkg.go.dev/badge/github.com/simplylib/gitsum.svg)](https://pkg.go.dev/github.com/simplylib/gitsum)
[![Go Report Card](https://goreportcard.com/badge/github.com/simplylib/gitsum)](https://goreportcard.com/report/github.com/simplylib/gitsum)

gitsum is a cli tool for finding information about large number of git repositories in one shot, fairly quickly.

## Installing

### Requirements
```git```


```
go install github.com/simplylib/gitsum@latest
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
