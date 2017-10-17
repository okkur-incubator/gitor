package main

import (
	"flag"
	"fmt"

	"github.com/okkur/gitor/commands"
)

func main() {
	//Flags that are to be added to commands
	var (
		pull string
	)

	var (
		remote string
		branch string
	)

	commands.Pull()

	flag.StringVar(&pull, "command", "gitor --command=pull", "pulls a repository")
	flag.StringVar(&remote, "remote", "gitor", "specifies remotename")
	flag.StringVar(&branch, "branch", "master", "specifies branchname")

	flag.Parse()

	fmt.Println(pull)
	fmt.Println(remote)
	fmt.Println(branch)
}
