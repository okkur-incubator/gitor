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
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	// "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const upstreamDefaultRemoteName string = "upstream"
const downstreamDefaultRemoteName string = "downstream"
const headRefBase string = "refs/heads/"
const remoteRefBase string = "refs/remotes/"

func main() {

	var (
		upstream      string
		upstreamRef   string
		downstreamRef string
		username      string
		token         string
		downstream    string
		localPath     string
	)

	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&upstreamRef, "upstreamRef", "master", "specifies upstream branch")
	flag.StringVar(&downstreamRef, "downstreamRef", "master", "specifies downstream branch")
	flag.StringVar(&username, "username", "", "specifies username")
	flag.StringVar(&token, "token", "", "specifies token or password")
	flag.StringVar(&downstream, "downstream", "", "specifies downstream")
	flag.StringVar(&localPath, "localPath", "", "specifies local repo filepath")
	flag.Usage = usage

	flag.Parse()

	command := flag.Arg(0)
	switch {
	case command == "update":
		username, token = checkEnvs(username, token)
		upstreamAuth := authType(upstream, username, token)
		downstreamAuth := authType(downstream, username, token)
		update(upstream, upstreamRef, downstream, downstreamRef, upstreamAuth, downstreamAuth, localPath)
	default:
		usage()
	}

}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gitor [subcommand] [flags]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Subcommands: update\n")
	os.Exit(2)
}

func checkEnvs(username string, token string) (string, string) {
	userEnv := os.Getenv("GITOR_USER")
	if username == "" {
		if userEnv == "" {
			log.Fatal("username not set")
		}
		username = userEnv
	}

	tokenEnv := os.Getenv("GITOR_TOKEN")
	if token == "" {
		if tokenEnv == "" {
			log.Fatal("token or password not set")
		}
		token = tokenEnv
	}
	return username, token
}

func authType(repo string, username string, token string) transport.AuthMethod {
	var auth transport.AuthMethod
	endpoint, err := transport.NewEndpoint(repo)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case endpoint.Protocol == "ssh":
		user := os.Getenv("USER")
		auth, err = ssh.NewSSHAgentAuth(user)
		if err != nil {
			log.Fatal(err)
		}
/* 	case endpoint.Protocol == "https":
		auth = http.NewBasicAuth(username, token) */
	default:
		auth = nil
	}

	return auth
}
