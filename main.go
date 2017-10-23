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
	"os"
)

func main() {

	var (
		upstream string
		branch   string
	)

	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&branch, "branch", "master", "specifies branch")
	flag.Usage = usage

	flag.Parse()
	command := flag.Args()[0]
	if command == "update" {
		update(upstream, branch)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gitor [subcommand] [flags]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Subcommands: update\n")
	os.Exit(2)
}
