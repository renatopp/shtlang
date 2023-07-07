package lang

import "sht/lang/runtime"

func Eval(input []byte) (string, error) {
	tree, err := Parse(input)
	if err != nil {
		return "", err
	}

	runtime := runtime.CreateRuntime()
	res, err := runtime.Run(tree)
	if err != nil {
		return "", err
	}

	return res, nil
}
