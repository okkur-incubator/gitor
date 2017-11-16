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
)

func main() {

	var (
		upstream   string
		branch     string
		username   string
		token      string
		downstream string
	)

	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&branch, "branch", "master", "specifies branch")
	flag.StringVar(&username, "username", username, "specifies username")
	flag.StringVar(&token, "token", token, "specifies token or password")
	flag.StringVar(&downstream, "downstream", downstream, "specifies downstream")
	flag.Usage = usage

	flag.Parse()

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

	command := flag.Arg(0)
	switch {
	case command == "update":
		update(upstream, branch, username, token, downstream)
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
