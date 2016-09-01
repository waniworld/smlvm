package parse

import (
	"e8vm.io/e8vm/asm8/ast"
	"e8vm.io/e8vm/lexing"
)

func parseVarStmts(p *parser, v *ast.Var) {
	for !(p.See(Rbrace) || p.See(lexing.EOF)) {
		stmt := parseVarStmt(p)
		if stmt != nil {
			v.Stmts = append(v.Stmts, stmt)
		}
		p.BailOut()
	}
}

func parseVar(p *parser) *ast.Var {
	ret := new(ast.Var)

	ret.Kw = p.ExpectKeyword("var")
	ret.Name = p.Expect(Operand)

	if ret.Name != nil {
		name := ret.Name.Lit
		if !IsIdent(name) {
			p.Errorf(ret.Name.Pos, "invalid var name %q", name)
		}
	}

	ret.Lbrace = p.Expect(Lbrace)
	if p.skipErrStmt() {
		return ret
	}

	parseVarStmts(p, ret)

	ret.Rbrace = p.Expect(Rbrace)
	ret.Semi = p.Expect(Semi)
	p.skipErrStmt()

	return ret
}
