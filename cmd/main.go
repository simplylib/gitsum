// Package cmd provides the cli functionality of gitsum.
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/simplylib/gitsum/git"
)

// Main starts the cli of gitsum.
func Main() error {
	log.Default().SetFlags(0)

	verbose := flag.Bool("v", false, "be verbose")
	help := flag.Bool("help", false, "show help message")
	flag.CommandLine.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: "+os.Args[0]+" <flags> [path]")
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	if *help {
		flag.PrintDefaults()
		return nil
	}

	var repoPath string

	if args := flag.Args(); len(args) != 1 {
		var err error
		repoPath, err = os.Getwd()

		if err != nil {
			return fmt.Errorf("could not get current directory (%w)", err)
		}
	} else {
		repoPath = filepath.Clean(args[0])
	}

	stat, err := os.Stat(repoPath)
	if err != nil {
		return fmt.Errorf("could not read path (%w)", err)
	}

	if !stat.IsDir() {
		return errors.New("path is not a directory")
	}

	repos, err := git.WalkDirForModifiedRepos(repoPath, *verbose)
	if err != nil {
		return fmt.Errorf("error during WalkDir (%w)", err)
	}

	if len(repos) == 0 {
		return nil
	}

	log.Println("Repos with changes:")

	for i := range repos {
		log.Println(repos[i])
	}

	return nil
}
