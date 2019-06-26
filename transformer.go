package shrinker

type Transformer interface {
	Transform(input []string) ([]string, error)
}

type TildeTransformer struct{}

func (tt *TildeTransformer) Transform(input []string) ([]string, error) {
	return nil, nil
}
