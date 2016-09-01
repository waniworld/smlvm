package parse

import (
	"io"

	"e8vm.io/e8vm/lexing"
)

// Tokens parses a file into a token array
func Tokens(f string, r io.Reader) ([]*lexing.Token, []*lexing.Error) {
	x, _ := makeTokener(f, r, false)
	toks := lexing.TokenAll(x)
	if errs := x.Errs(); errs != nil {
		return nil, errs
	}
	return toks, nil
}
