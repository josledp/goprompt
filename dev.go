package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-billy.v3/osfs"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/filesystem"
)

func findGitRepository() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("error determining current pwd: %v", err)
	}
	p := strings.Split(pwd, "/")
	for i := range p {
		path := strings.Join(p[:len(p)-i], "/")
		if info, err := os.Stat(path + "/.git"); os.IsNotExist(err) {
			continue
		} else if err == nil && info.IsDir() {
			return path, nil
		}
	}
	return "", fmt.Errorf("unable to find .git directory")
}

func main() {
	gitPwd, err := findGitRepository()
	if err != nil {
		log.Fatalf("not git repository: %v", err)
	}
	fs := osfs.New(gitPwd + "/.git")
	storage, err := filesystem.NewStorage(fs)
	if err != nil {
		log.Fatalf("unable to get storer: %v", err)
	}
	repository, err := git.Open(storage, osfs.New(gitPwd))
	if err != nil {
		log.Fatalf("unable to open repository: %v", err)
	}
	wt, err := repository.Worktree()
	if err != nil {
		log.Printf("error getting worktree: %v", err)
	}
	st, err := wt.Status()
	if err != nil {
		log.Printf("error getting worktree status: %v", err)
	}
	stage := 0
	work := 0
	untracked := 0
	for _, status := range st {
		switch status.Staging {
		case git.Unmodified:
		case git.Untracked:
		default:
			stage++
		}

		switch status.Worktree {
		case git.Unmodified:
		case git.Untracked:
			untracked++
		default:
			work++
		}
	}
	fmt.Printf("stage: %d wt: %d un:%d\n", stage, work, untracked)
	//stashes count lines from .git/logs/refs/stash
}
