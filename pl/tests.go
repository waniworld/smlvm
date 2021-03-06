package pl

import (
	"strings"

	"shanhu.io/smlvm/pl/tast"
	"shanhu.io/smlvm/pl/types"
	"shanhu.io/smlvm/syms"
)

func isTestName(name string) bool {
	if len(name) < len("TestX") {
		return false
	}
	if !strings.HasPrefix(name, "Test") {
		return false
	}
	lead := name[4]
	if lead >= 'a' && lead <= 'z' {
		return false
	}
	return true
}

func listTests(tops *syms.Table) []*objFunc {
	var list []*objFunc

	syms := tops.List()
	for _, s := range syms {
		if s.Type != tast.SymFunc {
			continue
		}
		f := s.Obj.(*objFunc)
		if f.isMethod {
			panic("bug") // a top level function should never be a method
		}
		if !types.SameType(f.ref.Type(), types.VoidFunc) {
			continue
		}
		name := s.Name()
		if isTestName(name) {
			list = append(list, f)
		}
	}

	return list
}
