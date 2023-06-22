package lang

type Evaluator struct {
	parser *Parser
}

func CreateEvaluator() *Evaluator {
	e := &Evaluator{}
	e.parser = CreateParser()
	return e
}

func Eval(input []byte) (string, error) {
	e := CreateEvaluator()
	r, err := e.Eval(input)
	return r, err
}

func (e *Evaluator) Eval(input []byte) (string, error) {
	tree, err := e.parser.Parse(input)

	return tree.String(), err
}
