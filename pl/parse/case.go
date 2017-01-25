package parse

import (
	"shanhu.io/smlvm/lexing"
	"shanhu.io/smlvm/pl/ast"
)

func parseCases(p *parser) []*ast.Case {
	var ret []*ast.Case
	for !(p.SeeOp("}") || p.See(lexing.EOF)) {
		if c := parseCase(p); c != nil {
			ret = append(ret, c)
		}
		p.skipErrStmt()
	}
	return ret
}

func parseCase(p *parser) *ast.Case {
	ret := new(ast.Case)
	if p.SeeKeyword("case") {
		ret.Kw = p.Shift()
		ret.Expr = parseExpr(p)
		if ret.Expr == nil {
			return nil
		}
	} else if p.SeeKeyword("default") {
		ret.Kw = p.Shift()
		ret.Expr = nil
	} else {
		p.CodeErrorfHere("pl.missingCaseInSwitch",
			"must start with keyword case/default in switch")
	}
	ret.Colon = p.ExpectOp(":")
	if ret.Colon == nil {
		return ret
	}
	for !(p.SeeKeyword("case") || p.SeeKeyword("default") ||
		p.SeeOp("}") || p.See(lexing.EOF)) {
		if stmt := p.parseStmt(); stmt != nil {
			ret.Stmts = append(ret.Stmts, stmt)
		}
		p.skipErrStmt()
	}
	return ret
}
