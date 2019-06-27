package main

import (
	"testing"

	shrinker "github.com/oinume/path-shrinker"
)

func TestRun(t *testing.T) {
	tests := map[string]struct {
		path   string
		config *shrinker.Config
		want   string
	}{
		"simple": {
			path: "/home/oinume/go",
			config: &shrinker.Config{
				Mode: shrinker.ModeShort,
			},
			want: "/h/o/g",
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
