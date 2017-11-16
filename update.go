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

func update(upstream string, branch string, username string, token string) error {

	// Validate URL
	err := validateUpstream(upstream, username, token)
	if err != nil {
		log.Println(err)
	}

	path := extractPath(upstream)

	// Initialize non bare repo
	r, err := git.PlainInit(path, false)
	if err != git.ErrRepositoryAlreadyExists && err != nil {
		log.Fatal(err)
	}

	// Open repo, if initialized
	r, err = git.PlainOpen(path)
	if err != nil {
		log.Println(err)
	}

	// Add a new remote, with the default fetch refspec
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{upstream},
	})

	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	// Pull using default options
	// If authentication required pull using authentication
	// TODO: needs switch for https:basicauth and ssh:keyauth
	err = w.Pull(&git.PullOptions{})
	if err != nil {
		switch err {
		case transport.ErrAuthenticationRequired:
			auth := http.NewBasicAuth(username, token)
			err = w.Pull(&git.PullOptions{Auth: auth})
			if err != nil {
				log.Fatal(err)
			}
		case transport.ErrEmptyRemoteRepository:
			log.Fatal("upstream repository is empty")
		default:
			log.Println(err)
		}
	}

	// Print the latest commit that was just pulled
	// TODO: simplify to only print commit hash "pulled: $hash"
	ref, err := r.Head()
	if err != nil {
		log.Println(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Println(err)
	}

	fmt.Printf("pulled: %s\n", commit.Hash)

	// Push using default options
	// If authentication required push using authentication
	// TODO: needs switch for https:basicauth and ssh:keyauth
	err = r.Push(&git.PushOptions{})
	if err != nil {
		switch err {
		case transport.ErrAuthenticationRequired:
			auth := http.NewBasicAuth(username, token)
			err = r.Push(&git.PushOptions{Auth: auth})
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(err)
		}
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

	return filePath
}

func validateUpstream(upstream string, username string, token string) error {
	// Create a temporary repository
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return err
	}

	// Add a new remote, with the default fetch refspec
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{upstream},
	})
	if err != nil {
		return err
	}

	// Fetch using the new remote
	// With authentication error use authentication
	// TODO: needs switch for https:basicauth and ssh:keyauth
	err = r.Fetch(&git.FetchOptions{})
	if err != nil {
		switch err {
		case transport.ErrAuthenticationRequired:
			auth := http.NewBasicAuth(username, token)
			err = r.Fetch(&git.FetchOptions{
				RemoteName: git.DefaultRemoteName,
				Auth:       auth,
			})
			if err != nil {
				return err
			}
		default:
			return err
		}
	}

	return nil
}
