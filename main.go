package main

import (
	"flag"
)

func main() {

	var (
		scmd     string
		upstream string
		branch   string
	)

	flag.StringVar(&scmd, "command", "pull", "pulls a repository")
	flag.StringVar(&upstream, "upstream", "https://github.com/okkur/gitor.git", "specifies upstream")
	flag.StringVar(&branch, "branch", "master", "specifies branch")

	flag.Parse()

	if scmd == "pull" {
		pull(scmd, upstream, branch)
	}
}
