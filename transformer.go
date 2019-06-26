package shrinker

type Transformer interface {
	Transform(input []string) []string
}

type TildeTransformer struct{}

func (tt *TildeTransformer) Transform(input []string) []string {
	return nil
}
