package main

import (
	"fmt"
	"os"

	"github.com/piokaczm/tools/siu/docker"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("error: not enough arguments")
		os.Exit(1)
	}

	t, err := docker.New(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	v, err := t.Translate()
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(v)
	os.Exit(0)
}
