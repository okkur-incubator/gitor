package main

import (
	"fmt"
)

func pull(scmd string, upstream string, branch string) error {
	fmt.Println(scmd)
	fmt.Println(upstream)
	fmt.Println(branch)
	return nil
}
