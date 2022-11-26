package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil" //nolint:staticcheck
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"

	shrinker "github.com/oinume/path-shrinker"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
	builtBy = ""
)

const (
	ExitOK    = 0
	ExitError = 1
)

type cli struct {
	outStream   io.Writer
	errStream   io.Writer
	readDirFunc shrinker.ReadDirFunc
}

func newCLI(outStream, errStream io.Writer, readDirFunc shrinker.ReadDirFunc) *cli {
	return &cli{
		outStream:   outStream,
		errStream:   errStream,
		readDirFunc: readDirFunc,
	}
}

func (c *cli) run(args []string) int {
	flagSet := flag.NewFlagSet("path-shrinker", flag.ContinueOnError)
	flagSet.SetOutput(c.errStream)
	var (
		fish         = flagSet.Bool("fish", false, "Enable -short -tilde -last")
		last         = flagSet.Bool("last", false, "Print the last directory's full name.")
		short        = flagSet.Bool("short", false, "Truncate directory names to the first character. Without -short, names are truncated without making them ambiguous.")
		tilde        = flagSet.Bool("tilde", false, "Substitute ~ for the home directory.")
		printVersion = flagSet.Bool("version", false, "Print current version.")
	)

	if err := flagSet.Parse(args[1:]); err != nil {
		return ExitError
	}

	if *printVersion {
		_, _ = fmt.Fprintf(c.outStream, "path-shrinker\n%s\n", c.getVersion(version, commit, date, builtBy))
		return ExitOK
	}

	if *fish {
		*tilde = true
		*last = true
		*short = true
	}
	config := &shrinker.Config{}
	if *tilde {
		config.ReplaceTilde = true
	}
	if *last {
		config.PreserveLast = true
	}
	if *short {
		config.Mode = shrinker.ModeShort
	} else {
		config.Mode = shrinker.ModeAmbiguous
	}

	var path string
	if len(flagSet.Args()) > 0 {
		path = flagSet.Args()[0]
	} else {
		p, err := os.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(c.errStream, "failed to get current working diretory: %v\n", err)
			return ExitError
		}
		path = p
	}

	result, err := c.shrinkPath(path, config)
	if err != nil {
		_, _ = fmt.Fprintf(c.errStream, "failed to run: %v\n", err)
		return ExitError
	}
	_, _ = fmt.Fprintln(c.outStream, result)

	return ExitOK
}

func main() {
	os.Exit(newCLI(
		os.Stdout,
		os.Stderr,
		ioutil.ReadDir,
	).run(os.Args))
}

func (c *cli) shrinkPath(path string, config *shrinker.Config) (string, error) {
	dirs := strings.Split(path, string(os.PathSeparator))
	transformers, err := c.createTransformers(dirs, config)
	if err != nil {
		return "", err
	}
	shrink, err := c.executeTransform(transformers, dirs, config)
	if err != nil {
		return "", err
	}
	return shrink, nil
}

func (c *cli) createTransformers(dirs []string, config *shrinker.Config) ([]shrinker.Transformer, error) {
	// -tilde, -short, -last are enabled
	// -> Process order: tilde, short, last
	// -amb, -last are enabled
	// -last just override last element with original value
	homeDir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}

	transformers := make([]shrinker.Transformer, 0, 4)
	if config.ReplaceTilde {
		transformers = append(transformers, &shrinker.ReplaceTildeTransformer{
			HomeDir: homeDir,
		})
	}

	switch config.Mode {
	case shrinker.ModeAmbiguous:
		var startDir string
		if config.ReplaceTilde {
			startDir = homeDir
		} else {
			startDir = "/" // TODO: Windows
		}
		transformers = append(transformers, &shrinker.AmbiguousTransformer{
			StartDir:    startDir,
			ReadDirFunc: c.readDirFunc,
		})
	case shrinker.ModeShort:
		transformers = append(transformers, &shrinker.ShortenTransformer{})
	default:
		return nil, fmt.Errorf("unknown mode: %v", config.Mode)
	}

	if config.PreserveLast {
		transformers = append(transformers, &shrinker.PreserveLastTransformer{
			Last: dirs[len(dirs)-1],
		})
	}

	return transformers, nil
}

func (c *cli) executeTransform(transformers []shrinker.Transformer, input []string, config *shrinker.Config) (string, error) {
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
	if !config.ReplaceTilde {
		path = sep + path
	}
	return path, nil
}

func (c *cli) getVersion(version, commit, date, builtBy string) string {
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
