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

	"gopkg.in/src-d/go-git.v4/plumbing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"

	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

func update(upstream string, upstreamRef string, downstream string, downstreamRef string,
	upstreamAuth transport.AuthMethod, downstreamAuth transport.AuthMethod, localPath string) error {

	// Validate upstream URL
	err := validateRepo(upstream, upstreamDefaultRemoteName, upstreamAuth)
	if err != nil {
		log.Println(err)
	}

	// Initialize non bare repo
	r, err := git.PlainInit(upstream, false)
	if err != git.ErrRepositoryAlreadyExists && err != nil {
		log.Fatal(err)
	}

	// Open repo, if initialized
	r, err = git.PlainOpen(upstream)
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

	push(r, downstream, upstreamRef, downstreamRef, downstreamAuth, localPath)

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

	reference := plumbing.ReferenceName(fmt.Sprintf("%s%s", headRefBase, upstreamRef))
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

func push(r *git.Repository, downstream string, upstreamRef string,
	downstreamRef string, downstreamAuth transport.AuthMethod, localPath string) {
	// Validate downstream URL
	err := validateRepo(downstream, downstreamDefaultRemoteName, downstreamAuth)
	if err != nil {
		log.Println(err)
	}

	// Push using default options
	// If authentication required push using authentication
	b := checkReference(r, downstream, downstreamRef, localPath, downstreamAuth)
	referenceList := []config.RefSpec{config.RefSpec(b+":"+b)}

	log.Printf("Pushing to %s ...\n", downstream)
	err = r.Push(&git.PushOptions{
		RemoteName: downstreamDefaultRemoteName,
		RefSpecs:   referenceList,
		Auth:       downstreamAuth,
	})
	if err != nil {
		log.Fatal(err)
	}
	remote, err := r.Reference(plumbing.ReferenceName(fmt.Sprintf("%s%s", headRefBase, downstreamRef)), true)
	if err != nil {
		log.Println(err)
	}
	remoteHash := remote.Hash()
	log.Printf("Pushed hash: %s\n", remoteHash)
	log.Println("Repository successfully synced")
}

func validateRepo(repo string, remoteName string, auth transport.AuthMethod) error {
	// Create a temporary repository
	r, err := git.Init(memory.NewStorage(), nil)
	if err != nil {
		return err
	}
	
	// Add a new remote, with the default fetch refspec
	_, err = r.CreateRemote(&config.RemoteConfig{
		Name: remoteName,
		URLs: []string{repo},
	})
	if err != nil {
		return err
	}
	
	// Fetch using the new remote
	err = r.Fetch(&git.FetchOptions{
		RemoteName: remoteName,
		Auth:       auth,
	})
	if err != nil {
		return err
	}

	return nil
}

func checkReference(r *git.Repository, downstream string, downstreamRef string,
	 localPath string, auth transport.AuthMethod) plumbing.ReferenceName {
	// Open an existing repository in a specific folder
	r, err := git.PlainOpen(localPath)
	if err != nil {
		log.Println(err)
	}
	
	ds, err := r.Remote(downstreamDefaultRemoteName)
	if err == git.ErrRemoteNotFound {
		ds, err = r.CreateRemote(&config.RemoteConfig{
			Name: downstreamDefaultRemoteName,
			URLs: []string{downstream},
		})
	}
	if err != nil {
		log.Println(err)
	}

	b := plumbing.ReferenceName(fmt.Sprintf("%s%s", headRefBase, downstreamRef))
	
	// Check if reference exists locally
	refs, err := r.References()
	if err != nil {
		log.Println(err)
	}
	var foundLocal bool

	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name() == b {
			log.Printf("reference exists locally:\n%s\n", ref)
			foundLocal = true
		}
		return nil
	})
	if !foundLocal {
		log.Printf("reference %s does not exist locally\n", b)
	}
	
	// Check if reference exists on remote
	remoteRefs, err := ds.List(&git.ListOptions{Auth: auth,})
	if err != nil {
		log.Println(err)
	}
	var found bool
	for _, ref := range remoteRefs {
		if ref.Name() == b {
			log.Printf("reference already exists on remote:\n%s\n", ref)
			found = true
		}
	}
	if !found {
		log.Printf("reference %s does not exist on remote\n", b)
	}

	return b
}
