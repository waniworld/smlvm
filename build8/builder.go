package build8

import (
	"e8vm.io/e8vm/debug8"
	"e8vm.io/e8vm/lex8"
	"e8vm.io/e8vm/link8"
)

// Builder builds a bunch of packages.
type Builder struct {
	*context
}

// NewBuilder creates a new builder with a particular home directory
func NewBuilder(input Input, output Output) *Builder {
	return &Builder{
		context: &context{
			input:      input,
			output:     output,
			pkgs:       make(map[string]*pkg),
			deps:       make(map[string][]string),
			linkPkgs:   make(map[string]*link8.Pkg),
			debugFuncs: debug8.NewFuncs(),
			Options:    new(Options),
		},
	}
}

// BuildPkgs builds a list of packages
func (b *Builder) BuildPkgs(pkgs []string) []*lex8.Error {
	return build(b.context, pkgs)
}

// Build builds a package.
func (b *Builder) Build(p string) []*lex8.Error {
	return b.BuildPkgs([]string{p})
}

// BuildPrefix builds packages with a particular prefix.
// in the path.
func (b *Builder) BuildPrefix(repo string) []*lex8.Error {
	return b.BuildPkgs(b.input.Pkgs(repo))
}

// BuildAll builds all packages.
func (b *Builder) BuildAll() []*lex8.Error { return b.BuildPrefix("") }

// LoadedPkgs is a loaded set of packages that is ready for building.
type LoadedPkgs struct {
	Targets []string
}

// Load loads a set of pckages that is ready for building.
func (b *Builder) Load(pkgs []string) *LoadedPkgs {
	panic("todo")
}
