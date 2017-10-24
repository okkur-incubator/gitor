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
	"net/url"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

func update(upstream string, branch string) error {
	path := parseURL(upstream)

	// We instance a new repository targeting the given path (the .git folder)
	r, err := git.PlainInit(path, false)
	if err != nil {
		log.Println(err)
	}

	fileSystemPath := parseFilesystemPath(path)
	r, err = git.PlainOpen(fileSystemPath)
	if err != nil {
		log.Println(err)
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
	if err != nil {
		log.Println(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Println(err)
	}

	fmt.Println(commit)
	return nil
}

func parseURL(upstream string) string {
	if strings.Contains(upstream, "git@github.com:") {
		url := "github.com/"
		result := strings.TrimPrefix(upstream, "git@github.com:")
		return url + result
	} else {
		u, err := url.Parse(upstream)
		if err != nil {
			log.Fatal(err)
		}

		result := strings.TrimSuffix(u.Host+u.Path, ".git")
		return result
	}
}

func parseFilesystemPath(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		log.Fatal(err)
	}

	result := u.Host + u.Path
	return result
}
