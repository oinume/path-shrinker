package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	path_shrinker "github.com/oinume/path-shrinker"
	shrinker "github.com/oinume/path-shrinker"
)

// https://github.com/robbyrussell/oh-my-zsh/tree/master/plugins/shrink-path
/*
$ pwd
/Users/kazuhiro/go/src/github.com/oinume/path-shrinker
$ shrink_path
/Use/k/g/s/gi/oi/pa
 */
var (
	tilde = flag.Bool("tilde", false, " Substitute ~ for the home directory.")
)

func main() {
	flag.Parse()
	path, err := run(flag.Args())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err=%v\n", err)
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

	transformers := make([]shrinker.Transformer, 0, 10)
	if *tilde {
		transformers = append(transformers, &path_shrinker.TildeTransformer{
			HomeDir: os.Getenv("HOME"), // TODO: go-homedir
		})
	}
	dirs := strings.Split(path, string(os.PathSeparator))
	//fmt.Printf("dirs = %+v\n", dirs)

	shrink, err := executeTransform(transformers, dirs);
	if err != nil {
		return "", err
	}
	return shrink, nil
}

func executeTransform(transformers []shrinker.Transformer, input []string) (string, error) {
	result := input
	for _, t := range transformers {
		output, err := t.Transform(result)
		if err != nil {
			return "", err
		}
		result = output
	}
	return strings.Join(result, string(os.PathSeparator)), nil
}
