package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestCLI_Run(t *testing.T) {
	if err := os.Setenv("HOME", "/home/oinume"); err != nil {
		t.Fatalf("failed to Setenv: %v", err)
	}

	tests := map[string]struct {
		args           []string
		wantOutput     string
		wantExitStatus int
	}{
		"short": {
			args:           []string{"main", "-short", "/home/oinume/go"},
			wantOutput:     "/h/o/g",
			wantExitStatus: ExitOK,
		},
		"short tilde": {
			args:           []string{"main", "-short", "-tilde", "/home/oinume/go"},
			wantOutput:     "~/g",
			wantExitStatus: ExitOK,
		},
		"last short tilde": {
			args:           []string{"main", "-last", "-short", "-tilde", "/home/oinume/go/src"},
			wantOutput:     "~/g/src",
			wantExitStatus: ExitOK,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			bout := new(bytes.Buffer)
			berr := new(bytes.Buffer)
			c := newCLI(bout, berr)
			exitStatus := c.run(test.args)
			if got, want := exitStatus, test.wantExitStatus; got != want {
				t.Fatalf("cli.run returns unexpected exit status: got=%v, want=%v", got, want)
			}
			if got, want := strings.TrimSpace(bout.String()), test.wantOutput; got != want {
				t.Errorf("cli.run outputs unexpected text: got=%q, want=%q", got, want)
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
