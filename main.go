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
	if command == "pull" {
		pull(command, upstream, branch)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gitor [subcommand] [flags]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "Subcommands: pull\n")
	os.Exit(2)
}
