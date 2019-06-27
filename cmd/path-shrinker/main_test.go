package main

import (
	"os"
	"testing"

	shrinker "github.com/oinume/path-shrinker"
)

func TestRun(t *testing.T) {
	os.Setenv("HOME", "/home/oinume")

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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := run([]string{test.path}, test.config)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.want {
				t.Errorf("got %q but want %q", got, test.want)
			}
		})
	}
}
