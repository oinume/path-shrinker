package shrinker

import (
	"os"
	"strings"
)

type Transformer interface {
	Transform(input []string) ([]string, error)
}

type TildeTransformer struct {
	HomeDir string
}

func (tt *TildeTransformer) Transform(input []string) ([]string, error) {
	path := strings.Join(input, string(os.PathSeparator))
	path = strings.Replace(path, tt.HomeDir, "~", -1)
	return strings.Split(path, string(os.PathSeparator)), nil
}
