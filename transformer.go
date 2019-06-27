package shrinker

import (
	"os"
	"strings"
)

type Config struct {
	PreserveLast bool
	Shorten      bool
	ReplaceTilde bool
}

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

type ShortTransformer struct{}

func (st *ShortTransformer) Transform(input []string) ([]string, error) {
	length := len(input)
	result := make([]string, 0, length)
	for _, v := range input {
		result = append(result, string([]rune(v)[0]))
	}
	return result, nil
}

type PreserveLastTransformer struct {
	Last string
}

func (plt *PreserveLastTransformer) Transform(input []string) ([]string, error) {
	result := make([]string, len(input))
	copy(result, input)
	length := len(input)
	result[length-1] = plt.Last
	return result, nil
}
