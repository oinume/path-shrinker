package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()
	path, err := run(flag.Args())
	if err != nil {
		fmt.Fprintf(os.Stderr, "err=%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%v\n", path)
}

func run(args []string) (string, error) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		//path = os.Getenv("PWD")
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = p
	}

	dirs := filepath.SplitList(path)
	fmt.Printf("dirs = %v\n", dirs)
	return path, nil
}
