/*
Copyright 2017 - The Gitor Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
