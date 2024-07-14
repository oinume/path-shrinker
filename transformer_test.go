package shrinker

import (
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/oinume/path-shrinker/shrinker_test"
)

func TestReplaceTildeTransformer_Transform(t *testing.T) {
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
			tr := &ReplaceTildeTransformer{
				HomeDir: test.homeDir,
			}
			got, err := tr.Transform(test.input)
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

func TestShortenTransformer_Transform(t *testing.T) {
	tests := map[string]struct {
		input   []string
		want    []string
		wantErr error
	}{
		"normal": {
			input:   strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			want:    []string{"U", "o", "g", "s", "g"},
			wantErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tr := &ShortenTransformer{}
			got, err := tr.Transform(test.input)
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

func TestPreserveLastTransformer_Transform(t *testing.T) {
	tests := map[string]struct {
		last    string
		input   []string
		want    []string
		wantErr error
	}{
		"normal": {
			last:    "github.com",
			input:   []string{"U", "o", "g", "s", "g"},
			want:    []string{"U", "o", "g", "s", "github.com"},
			wantErr: nil,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tr := &PreserveLastTransformer{
				Last: test.last,
			}
			got, err := tr.Transform(test.input)
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

func TestAmbiguousTransformer_Transform(t *testing.T) {
	tests := map[string]struct {
		input       []string
		want        []string
		wantErr     error
		readDirFunc ReadDirFunc
	}{
		"similar directories exist": {
			input:   strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			want:    []string{"U", "o", "go", "s", "gith"},
			wantErr: nil,
			readDirFunc: func(dirname string) (infos []os.DirEntry, e error) {
				return []os.DirEntry{
					fs.FileInfoToDirEntry(shrinker_test.NewMockFileInfo("go2", 0, 0755, time.Now(), true)),
					fs.FileInfoToDirEntry(shrinker_test.NewMockFileInfo("git", 0, 0755, time.Now(), true)),
				}, nil
			},
		},
		"no similar directories exist": {
			input:   strings.Split("/Users/oinume/go/src/github.com", string(os.PathSeparator)),
			want:    []string{"U", "o", "g", "s", "g"},
			wantErr: nil,
			readDirFunc: func(dirname string) (entries []os.DirEntry, e error) {
				return nil, nil
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tr := &AmbiguousTransformer{
				StartDir:    "/",
				ReadDirFunc: test.readDirFunc,
			}
			got, err := tr.Transform(test.input)
			if test.wantErr == nil {
				if !reflect.DeepEqual(got, test.want) {
					t.Errorf("unexpected result: got=%+v but want=%+v", got, test.want)
				}
			} else {
				if !reflect.DeepEqual(err, test.wantErr) {
					t.Fatalf("unexpected error: got=%v, want=%v", err, test.wantErr)
				}
			}
		})
	}
}

func TestAmbiguousTransformer_getAmbiguousName(t *testing.T) {
	type fields struct {
		startDir    string
		readDirFunc ReadDirFunc
	}
	type args struct {
		parent string
		target string
	}
	tests := map[string]struct {
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		"similar_name_exists": {
			fields: fields{
				startDir: "a",
				readDirFunc: func(dirname string) ([]os.DirEntry, error) {
					info := shrinker_test.NewMockFileInfo("Usr", 0, 0755, time.Now(), true)
					return []os.DirEntry{fs.FileInfoToDirEntry(info)}, nil
				},
			},
			args: args{
				parent: "a",
				target: "Users",
			},
			want: "Use",
		},
		"no_similar_name": {
			fields: fields{
				startDir: "a",
				readDirFunc: func(dirname string) ([]os.DirEntry, error) {
					info := shrinker_test.NewMockFileInfo("Home", 0, 0755, time.Now(), true)
					return []os.DirEntry{fs.FileInfoToDirEntry(info)}, nil
				},
			},
			args: args{
				parent: "a",
				target: "Users",
			},
			want: "U",
		},
		"empty_target": {
			fields: fields{
				startDir: "a",
				readDirFunc: func(dirname string) ([]os.DirEntry, error) {
					info := shrinker_test.NewMockFileInfo("Home", 0, 0755, time.Now(), true)
					return []os.DirEntry{fs.FileInfoToDirEntry(info)}, nil
				},
			},
			args: args{
				parent: "a",
				target: "",
			},
			want: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			at := &AmbiguousTransformer{
				StartDir:    test.fields.startDir,
				ReadDirFunc: test.fields.readDirFunc,
			}
			got, err := at.getAmbiguousName(test.args.parent, test.args.target)
			if (err != nil) != test.wantErr {
				t.Errorf("AmbiguousTransformer.getAmbiguousName(): error=%v, wantErr=%v", err, test.wantErr)
				return
			}
			if got != test.want {
				t.Errorf("AmbiguousTransformer.getAmbiguousName(): got=%q, want=%q", got, test.want)
			}
		})
	}
}
