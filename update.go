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

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func update(upstream string, upstreamRef string, downstream string, downstreamRef string,
	upstreamAuth transport.AuthMethod, downstreamAuth transport.AuthMethod) error {

	// Validate upstream URL
	err := validateRepo(upstream, upstreamAuth)
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

	// Add a upstream remote
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: upstreamDefaultRemoteName,
		URLs: []string{upstream},
	})
	if err != nil {
		switch err {
		case git.ErrRemoteNotFound:
			log.Fatal("remote not found")
		default:
			log.Println(err)
		}
	}

	pull(r, upstream, upstreamRef, upstreamAuth)

	push(r, downstream, upstreamRef, downstreamRef, downstreamAuth)

	return nil
}

func pull(r *git.Repository, upstream string, upstreamRef string, upstreamAuth transport.AuthMethod) {
	// Get the working directory for the repository
	w, err := r.Worktree()
	if err != nil {
		log.Fatal(err)
	}

	// Pull using default options
	// If authentication required pull using authentication
	log.Printf("Pulling %s ...\n", upstream)

	var reference plumbing.ReferenceName
	reference = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", upstreamRef))

	err = w.Pull(&git.PullOptions{
		RemoteName:    upstreamDefaultRemoteName,
		ReferenceName: reference,
		Auth:          upstreamAuth,
	})
	if err != nil {
		switch err {
		case transport.ErrEmptyRemoteRepository:
			log.Fatal("upstream repository is empty")
		default:
			log.Println(err)
		}
	}

	// Print the latest commit that was just pulled
	ref, err := r.Head()
	if err != nil {
		log.Println(err)
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		log.Println(err)
	}

	log.Printf("Pulled: %s\n", commit.Hash)
}

func push(r *git.Repository, downstream string, upstreamRef string, downstreamRef string, downstreamAuth transport.AuthMethod) {
	// Validate downstream URL
	err := validateRepo(downstream, downstreamAuth)
	if err != nil {
		log.Println(err)
	}

	// Add a downstream remote
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: downstreamDefaultRemoteName,
		URLs: []string{downstream},
	})
	if err != nil {
		switch err {
		case git.ErrRemoteNotFound:
			log.Fatal("remote not found")
		default:
			log.Println(err)
		}
	}

	// Push using default options
	// If authentication required push using authentication
	referenceList := append([]config.RefSpec{}, config.RefSpec(upstreamRef+":"+downstreamRef))
	log.Printf("Pushing to %s ...\n", downstream)
	err = r.Push(&git.PushOptions{
		RemoteName: downstreamDefaultRemoteName,
		RefSpecs:   referenceList,
		Auth:       downstreamAuth,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Repository successfully synced")
}

func extractPath(repo string) string {
	sshS := "ssh://"
	stSet := strings.Contains(repo, "https://")
	sSet := strings.Contains(repo, "http://")
	ssSet := strings.Contains(repo, "ssh://")
	cSet := strings.ContainsAny(repo, ":")

	// Handle ssh protocol - no protocol + colon suggests ssh
	if !stSet && !sSet && !ssSet && cSet {
		repo = sshS + repo
	}

	u, err := url.Parse(repo)
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

func validateRepo(repo string, upstreamAuth transport.AuthMethod) error {
	// Create a temporary repository
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return err
	}

	// Add a new remote, with the default fetch refspec
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: git.DefaultRemoteName,
		URLs: []string{repo},
	})
	if err != nil {
		return err
	}

	// Fetch using the new remote
	err = r.Fetch(&git.FetchOptions{
		RemoteName: git.DefaultRemoteName,
		Auth:       upstreamAuth,
	})
	if err != nil {
		return err
	}

	return nil
}
