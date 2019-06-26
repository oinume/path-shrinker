package shrinker

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestTildeTransformer_Transform(t *testing.T) {
	tests := map[string]struct {
		homeDir string
		input   []string
		want    []string
		wantErr error
	}{
		"tilde replacement": {
			homeDir: "/Users/oinume",
			input:   strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			want:    strings.Split("~/go/src/github.com", string(os.PathSeparator)),
			wantErr: nil,
		},
		"no tilde replacement": {
			homeDir: "/Users/oinuma",
			input:   strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			want:    strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			wantErr: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tt := &TildeTransformer{
				HomeDir: test.homeDir,
			}
			got, err := tt.Transform(test.input)
			if test.wantErr == nil {
				if !reflect.DeepEqual(got, test.want) {
					t.Errorf("got=%+v but want=%+v", got, test.want)
				}
			} else {
				if !reflect.DeepEqual(err, test.wantErr) {
					t.Fatal()
				}
			}
		})
	}
}
