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
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

const upstreamDefaultRemoteName string = "upstream"
const downstreamDefaultRemoteName string = "downstream"

func main() {

	var (
		upstream      string
		upstreamRef   string
		downstreamRef string
		user          string
		token         string
		downstream    string
		systemUser    string
		pemPath       string
		pemPwd        string
	)

	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&upstreamRef, "upstreamRef", "master", "specifies upstream branch")
	flag.StringVar(&downstreamRef, "downstreamRef", "master", "specifies downstream branch")
	flag.StringVar(&user, "user", "", "specifies GitHub/GitLab username")
	flag.StringVar(&token, "token", "", "specifies token or password")
	flag.StringVar(&downstream, "downstream", "", "specifies downstream")
	flag.StringVar(&systemUser, "systemUser", "", "specifies system username")
	flag.StringVar(&pemPath, "pemPath", "", "specifies path to pem file")
	flag.StringVar(&pemPwd, "pemPwd", "", "specifies pem file password")
	flag.Usage = usage

	flag.Parse()

	command := flag.Arg(0)
	switch {
	case command == "update":
		username, token, pemPath, pemPwd := checkEnvs(user, token, systemUser, pemPath, pemPwd)
		upstreamAuth := authType(upstream, username, token, pemPath, pemPwd)
		downstreamAuth := authType(downstream, username, token, pemPath, pemPwd)
		update(upstream, upstreamRef, downstream, downstreamRef, upstreamAuth, downstreamAuth)
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

func checkEnvs(username string, token string, systemUser string, pemPath string, pemPwd string) (string, string, string, string) {
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

	systemUserEnv := os.Getenv("USER")
	if systemUser == "" {
		if systemUserEnv == "" {
			log.Fatal("username not set")
		}
		username = systemUserEnv
	}

	pemEnv := os.Getenv("PEM_PATH")
	if pemPath == "" {
		if pemEnv == "" {
			log.Fatal("path to pem file not set")
		}
		username = userEnv
	}

	pemPwdEnv := os.Getenv("PEM_PWD")
	if pemPwd == "" {
		pemPwd = pemPwdEnv
	}

	return username, token, pemPath, pemPwd
}

func authType(repo string, username string, token string, pemPath string, pemPwd string) transport.AuthMethod {
	endpoint, err := transport.NewEndpoint(repo)
	if err != nil {
		log.Fatal(err)
	}
	switch {
	case endpoint.Protocol() == "ssh":
		switch {
		case pemPath != "":
			user := os.Getenv("USER")
			auth, err := ssh.NewPublicKeysFromFile(user, pemPath, pemPwd)
			if err != nil {
				log.Fatal(err)
			}
			return auth
		default:
			user := os.Getenv("USER")
			auth, err := ssh.NewSSHAgentAuth(user)
			if err != nil {
				log.Fatal(err)
			}
			return auth
		}
	case endpoint.Protocol() == "https":
		auth := http.NewBasicAuth(username, token)
		return auth
	default:
		var auth transport.AuthMethod
		auth = nil
		return auth
	}
}
