package shrinker

import (
	"fmt"
	"os"
	"path/filepath"
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

type AmbiguousTransformer struct {
	StartDir string
}

func (at *AmbiguousTransformer) getAmbiguousName(parent string, target string) (string, error) {
	result := target
	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		fmt.Printf("walk(): path=%v\n", path)
		if target == info.Name() {
			result = "[A]"
		}
		return nil
	}
	if err := filepath.Walk(target, walk); err != nil {
		return "", err
	}
	return result, nil
}

func (at *AmbiguousTransformer) Transform(input []string) ([]string, error) {
	//at.findAmbiguous
	for i, dir := range input {
		at.getAmbiguousName(at.StartDir, input[i])
	}

	fmt.Printf("input = %+v\n", input)
	//ioutil.ReadDir()
	walk := func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		fmt.Printf("walk(): path=%v\n", path)
		return nil
	}
	if err := filepath.Walk(at.StartDir, walk); err != nil {
		return nil, err
	}

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
