package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCLI_Run_OK(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		args       []string
		wantOutput string
	}{
		"short": {
			args:       []string{"main", "-short", "/home/oinume/go"},
			wantOutput: "/h/o/g",
		},
		"short tilde": {
			args:       []string{"main", "-short", "-tilde", "/home/oinume/go"},
			wantOutput: "~/g",
		},
		"last short tilde": {
			args:       []string{"main", "-last", "-short", "-tilde", "/home/oinume/go/src"},
			wantOutput: "~/g/src",
		},
		"fish": {
			args:       []string{"main", "-fish", "/home/oinume/go/src/github.com"},
			wantOutput: "~/g/s/github.com",
		},
		"tilde ambiguous": {
			args:       []string{"main", "-tilde", "/home/oinume/go/src/github.com"},
			wantOutput: "~/g/s/gi",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr)
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

func TestCLI_Run_FlagError(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		args       []string
		wantOutput string
	}{
		"unknown flag": {
			args:       []string{"main", "-unknown", "/home/oinume/go"},
			wantOutput: "/h/o/g",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr)
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

/*
func TestRun(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		path   string
		config *shrinker.Config
		want   string
	}{
		"short": {
			path: "/home/oinume/go",
			config: &shrinker.Config{
				Mode: shrinker.ModeShort,
			},
			want: "/h/o/g",
		},
		"tilde short": {
			path: "/home/oinume/go",
			config: &shrinker.Config{
				Mode:         shrinker.ModeShort,
				ReplaceTilde: true,
			},
			want: "~/g",
		},
		"tilde short last": {
			path: "/home/oinume/go/src",
			config: &shrinker.Config{
				Mode:         shrinker.ModeShort,
				ReplaceTilde: true,
				PreserveLast: true,
			},
			want: "~/g/src",
		},
		"tilde ambiguous": {
			path: "/home/oinume/go/src/github.com",
			config: &shrinker.Config{
				Mode:         shrinker.ModeAmbiguous,
				ReplaceTilde: true,
				PreserveLast: false,
			},
			want: "~/g/s/gi",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := run(test.path, test.config)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("got %q but want %q", got, test.want)
			}
		})
	}
}
*/
