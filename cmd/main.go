// Package cmd provides the cli functionality of gitsum.
package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type writtenToWriter bool

func (w *writtenToWriter) Write(p []byte) (int, error) {
	*w = true

	return len(p), nil
}

func isRepoModified(path string) (bool, error) {
	gitCmd := exec.Command("git", "-C", path, "status", "--porcelain")
	gitCmd.Stderr = io.Discard

	written := new(writtenToWriter)
	gitCmd.Stdout = written

	err := gitCmd.Run()
	if err != nil {
		return false, fmt.Errorf("git error (%w)", err)
	}

	if !*written {
		gitCmd = exec.Command("git", "-C", path, "log", "--branches", "--not", "--remotes")
		gitCmd.Stderr = io.Discard
		gitCmd.Stdout = written

		err := gitCmd.Run()
		if err != nil {
			return false, fmt.Errorf("git error checking branches for unpushed commits (%w)", err)
		}

		if *written {
			return true, nil
		}

		return false, nil
	}

	return true, nil
}

const gitNotARepoExitCode = 128

func isRepo(path string) (bool, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "--is-inside-work-tree")
	gitCmd.Stderr = io.Discard
	gitCmd.Stdout = io.Discard

	err := gitCmd.Run()
	if err != nil {
		var exitErr *exec.ExitError // todo: works on linux, unsure of others
		if !errors.As(err, &exitErr) {
			return false, fmt.Errorf("could not see if repo exists (%w)", err)
		}

		// exit code 128 is returned by git rev-parse --is-inside-worktree if current directory is not a
		if exitErr.ExitCode() == gitNotARepoExitCode {
			return false, nil
		}
	}

	return true, nil
}

// Main starts the cli of gitsum.
func Main() error {
	log.Default().SetFlags(0)

	verbose := flag.Bool("v", false, "be verbose")
	help := flag.Bool("help", false, "show help message")
	flag.CommandLine.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: gitsum <flags> [path]")
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

	var repos []string

	err = filepath.WalkDir(repoPath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			if *verbose {
				log.Printf("Could not read file (%v) error type (%t)\n", err, err)
			}
			return nil
		}

		if !entry.IsDir() {
			if *verbose {
				log.Printf("skipping (%v) as it is a file\n", path)
			}
			return nil
		}

		if *verbose {
			log.Printf("checking repo (%v)\n", path)
		}

		repo, err := isRepo(path)
		if err != nil {
			return fmt.Errorf("could not check if IsRepo (%w)", err)
		}

		if !repo {
			if *verbose {
				log.Printf("(%v) is not a repo", path)
			}
			return nil
		}

		modified, err := isRepoModified(path)
		if err != nil {
			return fmt.Errorf("could not check if repo was modified (%w)", err)
		}

		if !modified {
			if *verbose {
				log.Printf("(%v) is not modified\n", path)
			}
			return filepath.SkipDir
		}

		if *verbose {
			log.Printf("saw repo %v, skipping children\n", path)
		}

		repos = append(repos, entry.Name())

		return filepath.SkipDir
	})
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
