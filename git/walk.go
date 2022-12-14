package git

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type repoList struct {
	repos []string
	sync.Mutex
}

type walker struct {
	wg         *sync.WaitGroup
	list       *repoList
	filePath   string
	workingDir string
}

func (w *walker) walkFunc(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return fmt.Errorf("error during walkFunc (%w)", err)
	}

	if !entry.IsDir() {
		return nil
	}

	repo, err := IsRepo(path)
	if err != nil {
		return fmt.Errorf("could not check if IsRepo (%w)", err)
	}

	if !repo {
		if path == w.filePath {
			return nil
		}

		w.wg.Add(1)

		wkr := *w
		wkr.filePath = path

		go func(w walker) {
			defer w.wg.Done()

			if err := filepath.WalkDir(path, w.walkFunc); err != nil {
				log.Printf("could not walkdir (%v)", err)
			}
		}(wkr)

		return filepath.SkipDir
	}

	modified, err := IsRepoModified(path)
	if err != nil {
		return fmt.Errorf("could not check if repo was modified (%w)", err)
	}

	if modified {
		w.list.Lock()
		w.list.repos = append(w.list.repos, strings.TrimPrefix(path, w.workingDir+string(filepath.Separator)))
		w.list.Unlock()
	}

	return filepath.SkipDir
}

// WalkRepos returns a full path ex: (/home/<user>/git/<repo>) slice of modified repositories.
func WalkDirForModifiedRepos(filePath string, verbose bool) ([]string, error) {
	wg := &sync.WaitGroup{}
	repoList := &repoList{}

	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("could not get current directory (%w)", err)
	}

	wkr := walker{
		wg:         wg,
		list:       repoList,
		filePath:   filePath,
		workingDir: wd,
	}

	err = filepath.WalkDir(filePath, wkr.walkFunc)
	if err != nil {
		return nil, fmt.Errorf("could not filepath.WalkDir (%w)", err)
	}

	wg.Wait()

	return repoList.repos, nil
}
