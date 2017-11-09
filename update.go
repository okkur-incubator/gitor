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
	"strconv"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func update(upstream string, branch string, username string, password string) error {

	auth := http.NewBasicAuth(username, password)
	path := extractPath(upstream)

	// Create in memory repository, create remote and validate URL
	r, err := validateUpstream(upstream, auth, path)

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

	r, err = git.PlainOpen(path)
	if err != nil {
		log.Println(err)
	}

	// Push using default options or using authentication for https
	err = push(r)
	if err != nil {
		log.Fatal(err)
	}
	if err == transport.ErrAuthenticationRequired {
		err = r.Push(&git.PushOptions{Auth: auth})
	}

	return nil
}

func extractPath(upstream string) string {
	sshS := "ssh://"
	stSet := strings.Contains(upstream, "https://")
	sSet := strings.Contains(upstream, "http://")
	ssSet := strings.Contains(upstream, "ssh://")
	cSet := strings.ContainsAny(upstream, ":")

	// Handle ssh protocol - no protocol + colon suggests ssh
	if !stSet && !sSet && !ssSet && cSet {
		upstream = sshS + upstream
	}

	u, err := url.Parse(upstream)
	if err != nil {
		log.Fatal(err)
	}

	path := strings.TrimSuffix(u.Path, ".git")
	port := u.Port()

	// Handle parsing of part of path as port
	if _, err := strconv.Atoi(port); err != nil {
		path = port + path
	}

	// Check for missing path separator
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	host := u.Hostname()
	filePath := host + path

	fmt.Println(filePath)
	return filePath
}

func validateUpstream(upstream string, auth transport.AuthMethod, path string) (*git.Repository, error) {
	// Create a new repository
	r1, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		log.Println(err)
	}

	// Add a new remote, with the default fetch refspec
	_, err = r1.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{upstream},
	})

	// Fetch using the new remote
	err = r1.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})

	if err != nil {
		fmt.Println(err)
		log.Fatalf("%s is not a valid URL\n", upstream)
	}

	// We instance a new repository targeting the given path (the .git folder)
	r2, err := git.PlainInit(path, false)
	if err != nil {
		log.Println(err)
	}

	fileSystemPath := path
	r2, err = git.PlainOpen(fileSystemPath)
	if err != nil {
		log.Println(err)
	}

	// Add a new remote, with the default fetch refspec
	_, err = r2.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{upstream},
	})

	// Get the working directory for the repository
	w, err := r2.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	// Check if authentication required, pull the latest changes from the origin remote and merge into the current branch
	err = pull(w)
	if err != nil {
		log.Println(err)
	}
	if err == transport.ErrAuthenticationRequired {
		err = w.Pull(&git.PullOptions{Auth: auth})
	}

	return r2, nil
}

func pull(w *git.Worktree) error {
	err := w.Pull(&git.PullOptions{})
	if err != nil {
		return err
	}

	return nil
}

func push(r *git.Repository) error {

	// Push using default options
	err := r.Push(&git.PushOptions{})
	if err != nil {
		return err
	}

	return nil
}
