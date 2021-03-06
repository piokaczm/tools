package main

import (
	"fmt"
	"os"

	"github.com/piokaczm/tools/siu/docker/cli"
)

func main() {
	if err := cli.Invoke(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
