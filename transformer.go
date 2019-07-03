package shrinker

import (
	"os"
	"strings"
)

type Mode int

const (
	ModeAmbiguous Mode = iota + 1
	ModeShort
)

type Config struct {
	PreserveLast bool
	Shorten      bool
	ReplaceTilde bool
	Mode         Mode
}

type Transformer interface {
	Transform(input []string) ([]string, error)
}

type ReplaceTildeTransformer struct {
	HomeDir string
}

func (tt *ReplaceTildeTransformer) Transform(input []string) ([]string, error) {
	path := strings.Join(input, string(os.PathSeparator))
	path = strings.Replace(path, tt.HomeDir, "~", -1)
	return strings.Split(path, string(os.PathSeparator)), nil
}

type AmbiguousTransformer struct{}

func (at *AmbiguousTransformer) Transform(input []string) ([]string, error) {
	// TODO: check directory. https://flaviocopes.com/go-list-files/
	length := len(input)
	result := make([]string, 0, length)
	for i, v := range input {
		runes := []rune(v)
		if len(runes) == 0 {
			continue
		}
		if i == length-1 && len(runes) > 1 {
			result = append(result, string(runes[0])+string(runes[1]))
		} else {
			result = append(result, string(runes[0]))
		}
	}
	return result, nil
}

type ShortenTransformer struct{}

func (st *ShortenTransformer) Transform(input []string) ([]string, error) {
	length := len(input)
	result := make([]string, 0, length)
	for _, v := range input {
		runes := []rune(v)
		if len(runes) == 0 {
			continue
		}
		result = append(result, string(runes[0]))
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
