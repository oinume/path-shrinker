package shrinker

import (
	"math"
	"os"
	"path/filepath"
	"strings"
	"unicode"
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

type ReadDirFunc func(dirname string) ([]os.DirEntry, error)

type AmbiguousTransformer struct {
	StartDir    string
	ReadDirFunc ReadDirFunc
}

func (at *AmbiguousTransformer) getAmbiguousName(parent, target string) (string, error) {
	if target == "" {
		return "", nil
	}
	result := ""
	files, err := at.ReadDirFunc(parent)
	if err != nil {
		return "", err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		if f.Name() == target {
			continue
		}

		nameRunes := []rune(f.Name())
		targetRunes := []rune(target)
		maxLength := int(math.Max(float64(len(nameRunes)), float64(len(targetRunes))))
		minLength := int(math.Min(float64(len(nameRunes)), float64(len(targetRunes))))
		a := make([]rune, 0, maxLength)
		previousMatched := true
		var i int
		for i = 0; i < minLength; i++ {
			// Compare character case insensitively
			if previousMatched && unicode.ToLower(nameRunes[i]) == unicode.ToLower(targetRunes[i]) {
				a = append(a, targetRunes[i])
			} else {
				// TODO: maybe a bug
				previousMatched = false //nolint:ineffassign
				break
			}
		}
		if i < len(target)-1 {
			// Append additional 1 rune to distinguish name
			a = append(a, targetRunes[i])
		}

		if len(a) > len(result) {
			result = string(a)
		}
	}

	if len(result) == 0 {
		result = string([]rune(target)[0])
	}
	return result, nil
}

func (at *AmbiguousTransformer) Transform(input []string) ([]string, error) {
	parent := at.StartDir
	result := make([]string, 0, len(input))
	for _, dir := range input {
		if dir == "" {
			continue
		}
		// TODO: goroutine
		name, err := at.getAmbiguousName(parent, dir)
		if err != nil {
			return nil, err
		}
		result = append(result, name)
		parent = filepath.Join(parent, dir)
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
