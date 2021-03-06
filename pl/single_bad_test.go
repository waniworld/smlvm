package pl

import (
	"testing"

	"shanhu.io/smlvm/arch"
)

func TestSingleFileBad(t *testing.T) {
	o := func(code, input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		if es == nil {
			t.Log(input)
			t.Error("should error:", code)
			return
		}
		errNum := len(es)
		if errNum != 1 {
			t.Log(input)
			t.Logf("%d errors returned", errNum)
			for _, err := range es {
				t.Log(err.Code)
			}
		}
		code = "pl." + code
		if es[0].Code != code {
			t.Log(input)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
			return
		}
	}

	// test function for declConflict
	c := func(code, input string) {
		_, es, _ := CompileSingle("main.g", input, false)
		if es == nil || len(es) != 2 {
			t.Log(input)
			t.Error("should have 2 errors for confliction")
			return
		}
		code = "pl." + code
		if es[0].Code != code {
			t.Log(input)
			t.Errorf("ErrCode expected: %q, got %q", code, es[0].Code)
			return
		}

		const previousPos = "pl.declConflict.previousPos"
		if es[1].Code != previousPos {
			t.Log(input)
			t.Errorf("ErrCode expected: %q, got %q", previousPos, es[1].Code)
			return
		}
	}

	// no main
	o("missingMainFunc", "")

	// missing returns
	o("missingReturn", `func f() int { }`)
	o("missingReturn", `func f() int { for { break } }`)
	o("missingReturn", `func f() int { for { if true { break } } }`)
	o("missingReturn", `func f() int { for { if true break } }`)
	o("missingReturn",
		`func f() int { for true { if true { return 0 } } }`)
	o("missingReturn",
		`func f() int { for true { if true return 0 } }`)
	o("missingReturn", `func f() int { if true { return 0 } }`)
	o("missingReturn", `func f() int { if true return 0 }`)

	o("cannotAssign.typeMismatch", `func f() int {return 'a'}`)
	o("cannotAssign.typeMismatch", `func f() (int, int) {return 1, 'a'}`)
	o("return.expectNoReturn", `func f() {return 1}`)
	o("return.noReturnValue", `func f() int {return}`)

	// confliction return 2 errors
	c("declConflict.func", `func a() {}; func a() {}`)
	c("declConflict.field", `struct A { b int; b int }`)
	c("declConflict.const", `const a=1; const a=2`)
	c("declConflict.struct", "struct A{}; struct A{}")
	c("declConflict.interface", "interface A{}; interface A{}")
	c("declConflict.interfaceFunc", "interface A{ a(); a()}")
	c("declConflict.var", "var a int; func a() {}")
	c("declConflict.func", "func main() {}; func main() {};")

	// unused vars
	o("unusedSym", `func main() { var a int }`)
	o("unusedSym", `func main() { var a = 3 }`)
	o("unusedSym", `func main() { a := 3 }`)
	o("unusedSym", `func main() { var a int; a=3 }`)
	o("unusedSym", `func main() { var a int; (a)=(3) }`)
	o("unusedSym", `func main() { var a,b=3,4; _:=a }`)
	o("unusedSym", `func main() { for i:=1;;{} }`)

	// parser, import related
	o("multiImport", `import(); import()`)

	// expect ';', got keyword
	o("missingSemi", "import() func main(){}")

	// circular dependence
	o("circDep.struct", `struct A { a A };`)
	o("circDep.struct", `struct A { b B }; struct B { a A };`)
	o("circDep.struct", `struct A { b B }; struct B { a [3]A };`)
	o("circDep.struct", `struct A { b B }; struct B { a [0]A };`)
	o("circDep.const", `const a = b; const b = a`)
	o("circDep.const", `const a = 3 + b; const b = a`)
	o("circDep.const", `const a = 3 + b; const b = 0 - a`)
	o("circDep.const", "const a, b = a, b")

	// assign and allocate
	o("cannotAlloc", `struct A {}; func main() { a:=A }`)
	o("cannotAlloc", `struct A {}; func (a *A) f(){};
		func main() { var a A; f:=a.f; _:=f }`)

	o("cannotAssign.typeMismatch", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func()=a.f; _:= f }`)
	o("cannotAssign.typeMismatch", `struct A {}; func (a *A) f(){};
		func main() { var a A; var f func(); f=a.f; _:= f }`)
	o("cannotAssign.typeMismatch", `func main() { var a [2]int; var b [3]int;
		a=b}`)
	o("cannotAssign.interface",
		`interface I { t() int }
		 func main() {var i I; i=1; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t() int }
		 struct S { t int }
		 func main() {var i I; var s S; i=s; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t() int; }
		 struct S { i int }
		 func (s *S) t() {}
		 func main() {var i I; var s S; i=s; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t() int }
		 struct S { i int }
		 func (s *S) a() int {}
		 func main() {var i I; var s S; i=s; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t() int }
		 interface S { s() int }
		 func main() {var i I; var s S; i=s; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t() (int,int) }
		 interface S { t() int }
		 func main() {var i I; var s S; i=s; _:=i}`)
	o("cannotAssign.interface",
		`interface I { t(a int) }
		 interface S { t(a,b int ) }
		 func main() {var i I; var s S; i=s; _:=i}`)

	// others
	o("multiRefInExprList", ` func r() (int, int) { return 3, 4 }
		func p(a, b, c int) { }
		func main() { p(r(), 5) }`)

	o("invalidStructDecl", "type A struct {}")
	o("star.onNotSingle", `func main() { var a=*A() };
	func A() (*int, *int) { return nil, nil}`)
	o("incStmt.notSingle", ` func f() (int, int) { return 0, 0 }
		func main() { f()++ }`)

	// not single
	o("switchExpr.notSingle", ` func f() (int, int) { return 0, 0 }
		func main() { switch f(){case 1:} }`)
	o("caseExpr.notSingle", ` func f() (int, int) { return 0, 0 }
		func main() { a:=4; switch a {case f():} }`)

	// Bugs found by the fuzzer in the past
	o("undefinedIdent", "func f() **o.o {}")
	o("expectConstExpr", "func n()[char[:]]string{}")
	o("undefinedIdent", "const c, d = d, t; func main() {}")
	o("cannotCast", `func main() {
			var s string
			for i := 0; i < len(s-2); i++ {}
		}`)

	o("notYetSupported", `interface I { t() int }
		func main() { var b *B; var i I; i=b; _:=i;}
		struct B {
    		v int
		}
		func (b *B) t() int {
			return b.v
		}`)
	o("notYetSupported", `interface I { t() int }
		func main() { var b I2; var i I; i=b; _:=i;}
		interface I2 {t() int}`)

}

func TestSingleFilePanic(t *testing.T) {
	// runtime errors

	const N = 100000
	o := func(input string) {
		_, e := singleTestRun(t, input, N)
		if !arch.IsErr(e, arch.ErrPanic) {
			t.Log(input)
			t.Log(e)
			t.Error("should panic")
			return
		}
	}

	o("func main() { panic() }")
	o("func main() { var pa *int; printInt(*pa) }")
	o("struct A { a int }; func main() { var pa *A; b:=pa.a; _:=b }")
	o("func main() { var a func(); a() }")
	o("func f() {}; func main() { var a func()=f; a=nil; a() }")
	o("func f(p *int) { printInt(*p) }; func main() { f(nil) }")
	o("struct A { p *int }; func main() { var a A; a.p=nil; *a.p=0 }")
}
