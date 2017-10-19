package main

import (
	"fmt"
	"log"

	git "gopkg.in/src-d/go-git.v4"
)

func update(upstream string, branch string) error {

	path := upstream

	// We instance a new repository targeting the given path (the .git folder)
	r, err := git.PlainOpen(path)
	if err != nil {
		log.Fatal(err)
	}

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	// Pull the latest changes from the origin remote and merge into the current branch
	err = w.Pull(&git.PullOptions{RemoteName: upstream})

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())

	fmt.Println(commit)
	return nil
}
