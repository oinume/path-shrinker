package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

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
	short = flag.Bool("short", false, "Truncate directory names to the first character. Without -short, names are truncated without making them ambiguous.")
	tilde = flag.Bool("tilde", false, "Substitute ~ for the home directory.")
	last  = flag.Bool("last", false, "Print the last directory's full name.")
)

func main() {
	flag.Parse()
	path, err := run(flag.Args())
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err=%v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}

func run(args []string) (string, error) {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		p, err := os.Getwd()
		if err != nil {
			return "", err
		}
		path = p
	}

	dirs := strings.Split(path, string(os.PathSeparator))
	transformers := createTransformers(dirs)
	shrink, err := executeTransform(transformers, dirs)
	if err != nil {
		return "", err
	}
	return shrink, nil
}

func createTransformers(dirs []string) []shrinker.Transformer {
	// -tilde, -short, -last are enabled
	// -> Process order: tilde, short, last
	// -amb, -last are enabled
	// -last just override last element with original value
	transformers := make([]shrinker.Transformer, 0, 4)
	if *tilde {
		transformers = append(transformers, &shrinker.ReplaceTildeTransformer{
			HomeDir: os.Getenv("HOME"), // TODO: go-homedir
		})
	}
	if *short {
		transformers = append(transformers, &shrinker.ShortenTransformer{})
	}
	if *last {
		transformers = append(transformers, &shrinker.PreserveLastTransformer{
			Last: dirs[len(dirs)-1],
		})
	}
	return transformers
}

func executeTransform(transformers []shrinker.Transformer, input []string) (string, error) {
	result := input
	for _, t := range transformers {
		if len(result) == 0 {
			return "", fmt.Errorf("empty result was returned")
		}
		output, err := t.Transform(result)
		if err != nil {
			return "", err
		}
		result = output
	}
	return strings.Join(result, string(os.PathSeparator)), nil
}
