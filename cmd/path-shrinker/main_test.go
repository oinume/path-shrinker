package main

import (
	"bytes"
	"io/fs"
	"os"
	"strings"
	"testing"
	"time"

	shrinker "github.com/oinume/path-shrinker"
	"github.com/oinume/path-shrinker/shrinker_test"
)

func mockReadDir(dirname string) ([]os.DirEntry, error) {
	return nil, nil
}

func TestCLI_Run_OK(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		args        []string
		readDirFunc shrinker.ReadDirFunc
		wantOutput  string
	}{
		"short": {
			args:        []string{"main", "-short", "/home/oinume/go"},
			wantOutput:  "/h/o/g",
			readDirFunc: mockReadDir,
		},
		"short tilde": {
			args:        []string{"main", "-short", "-tilde", "/home/oinume/go"},
			wantOutput:  "~/g",
			readDirFunc: mockReadDir,
		},
		"last short tilde": {
			args:        []string{"main", "-last", "-short", "-tilde", "/home/oinume/go/src"},
			wantOutput:  "~/g/src",
			readDirFunc: mockReadDir,
		},
		"fish": {
			args:        []string{"main", "-fish", "/home/oinume/go/src/github.com"},
			wantOutput:  "~/g/s/github.com",
			readDirFunc: mockReadDir,
		},
		"tilde ambiguous": {
			args:        []string{"main", "-tilde", "/home/oinume/go/src/github.com"},
			wantOutput:  "~/g/s/g",
			readDirFunc: mockReadDir,
		},
		"ambiguous": {
			args:       []string{"main", "/home/oinume/go/src/github.com"},
			wantOutput: "/h/o/g/s/gith",
			readDirFunc: func(dirname string) ([]os.DirEntry, error) {
				fileInfo := shrinker_test.NewMockFileInfo("git", 0, 0755, time.Now(), true)
				ret := []os.DirEntry{
					fs.FileInfoToDirEntry(fileInfo),
				}
				return ret, nil
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr, test.readDirFunc)
			exitStatus := c.run(test.args)
			if got, want := exitStatus, ExitOK; got != want {
				t.Fatalf("cli.run returns unexpected exit status: got=%v, want=%v", got, want)
			}
			if got, want := strings.TrimSpace(bout.String()), test.wantOutput; got != want {
				t.Errorf("cli.run outputs unexpected text: got=%q, want=%q", got, want)
			}
			if got := berr.String(); got != "" {
				t.Errorf("cli.run outputs unexpected text to error stream: %q", got)
			}
		})
	}
}

func TestCLI_Run_PrintVersion(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	version = "1.0.0"
	commit = "abc"
	tests := map[string]struct {
		args       []string
		wantOutput string
	}{
		"version": {
			args:       []string{"main", "-version", "/home/oinume/go"},
			wantOutput: "path-shrinker\nversion: 1.0.0\ncommit: abc\n",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr, os.ReadDir)
			exitStatus := c.run(test.args)
			if got, want := exitStatus, ExitOK; got != want {
				t.Fatalf("cli.run returns unexpected exit status: got=%v, want=%v", got, want)
			}
			if got := bout.String(); got != test.wantOutput {
				t.Errorf("cli.run outputs unexpected text to stdout: got=%q, want=%q", got, test.wantOutput)
			}
			if berr.String() != "" {
				t.Errorf("cli.run outputs text to stderr somehow: %q", berr.String())
			}
		})
	}
}

func TestCLI_Run_FlagError(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		args []string
	}{
		"unknown flag": {
			args: []string{"main", "-unknown", "/home/oinume/go"},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr, os.ReadDir)
			exitStatus := c.run(test.args)
			if got, want := exitStatus, ExitError; got != want {
				t.Fatalf("cli.run returns unexpected exit status: got=%v, want=%v", got, want)
			}
			if got := bout.String(); got != "" {
				t.Errorf("cli.run outputs unexpected text to stdout: %q", got)
			}
			if berr.String() == "" {
				t.Errorf("cli.run does not output any text to stderr")
			}
		})
	}
}
