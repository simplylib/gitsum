package git

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
)

type writtenToWriter bool

func (w *writtenToWriter) Write(p []byte) (int, error) {
	*w = true

	return len(p), nil
}

func IsRepoModified(path string) (bool, error) {
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

func IsRepo(path string) (bool, error) {
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
