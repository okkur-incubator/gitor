package main

import (
	"fmt"
)

func pull(command string, upstream string, branch string) error {
	fmt.Println(command)
	fmt.Println(upstream)
	fmt.Println(branch)
	return nil
}
