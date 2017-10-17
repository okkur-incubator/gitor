package main

import (
	"flag"
)

func main() {

	var (
		upstream string
		branch   string
	)

	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&branch, "branch", "master", "specifies branch")

	flag.Parse()
	command := flag.Args()[0]
	if command == "pull" {
		pull(command, upstream, branch)
	}
}
