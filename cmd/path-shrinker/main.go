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
	short        = flag.Bool("short", false, "Truncate directory names to the first character. Without -short, names are truncated without making them ambiguous.")
	tilde        = flag.Bool("tilde", false, "Substitute ~ for the home directory.")
	last         = flag.Bool("last", false, "Print the last directory's full name.")
	printVersion = flag.Bool("version", false, "Print current version.")
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

func main() {
	flag.Parse()
	if *printVersion {
		fmt.Printf("path-shrinker\n%s\n", getVersion(version, commit, date, builtBy))
		os.Exit(0)
	}

	config := createConfig()
	path, err := run(flag.Args(), config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "err=%v\n", err)
		os.Exit(1)
	}
	fmt.Println(path)
}

func run(args []string, config *shrinker.Config) (string, error) {
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
	transformers := createTransformers(dirs, config)
	shrink, err := executeTransform(transformers, dirs)
	if err != nil {
		return "", err
	}
	return shrink, nil
}

func createTransformers(dirs []string, config *shrinker.Config) []shrinker.Transformer {
	// -tilde, -short, -last are enabled
	// -> Process order: tilde, short, last
	// -amb, -last are enabled
	// -last just override last element with original value
	transformers := make([]shrinker.Transformer, 0, 4)
	if config.ReplaceTilde {
		transformers = append(transformers, &shrinker.ReplaceTildeTransformer{
			HomeDir: os.Getenv("HOME"), // TODO: go-homedir
		})
	}

	switch config.Mode {
	case shrinker.ModeAmbiguous:
		transformers = append(transformers, &shrinker.ShortenTransformer{}) // TODO: ambiguous
	case shrinker.ModeShort:
		transformers = append(transformers, &shrinker.ShortenTransformer{})
	default:
		panic("Unknown mode")
	}

	if config.PreserveLast {
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
	const sep = string(os.PathSeparator)
	path := strings.Join(result, sep)
	if !strings.HasPrefix(path, "~") { // TODO: config
		path = sep + path
	}
	return path, nil
}

func getVersion(version, commit, date, builtBy string) string {
	var result = fmt.Sprintf("version: %s", version)
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	if builtBy != "" {
		result = fmt.Sprintf("%s\nbuilt by: %s", result, builtBy)
	}
	return result
}

func createConfig() *shrinker.Config {
	c := &shrinker.Config{}
	if *tilde {
		c.ReplaceTilde = true
	}
	if *last {
		c.PreserveLast = true
	}
	if *short {
		c.Mode = shrinker.ModeShort
	} else {
		c.Mode = shrinker.ModeAmbiguous
	}
	return c
}
