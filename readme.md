[![BuildStatus](https://travis-ci.org/e8vm/e8vm.png?branch=master)](https://travis-ci.org/e8vm/e8vm)

```
go get -u e8vm.io/e8vm/...
```

# E8VM

[![Join the chat at https://gitter.im/e8vm/e8vm](https://badges.gitter.im/e8vm/e8vm.svg)](https://gitter.im/e8vm/e8vm?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Emul8ed Virtual Machine (E8VM) is a self-contained system that has its
own instruction set -- `arch8`, its own assembly language and
assembler -- `asm8`, its own system language -- `g8`, and its own
project building system -- `build8`. Using `g8` and `build8`, we can
build a small operating system [`os8`](https://github.com/e8vm/os8).

The project is written entirely in Go language. Plus, each file in the
project has no more than 300 lines, with each line no more than 80
characters. Among these small files, there are no circular
dependencies, and as a result, the project architecture can be
automatically [visualized](http://8k.lonnie.io) from static code
analysis.

For Go language documentation on the package APIs, I recommend
[GoWalker](https://gowalker.org/e8vm.io/e8vm). I find it slightly
better than [godoc.org](https://godoc.org/e8vm.io/e8vm).

## To Use `make`

The project comes with a `makefile`, which formats the code files,
check lints, check circular dependencies and build tags. Running the
`makefile` requires installing some tools.

```
go get -u e8vm.io/tools/...
go get -u github.com/golang/lint/golint
go get -u github.com/jstemmer/gotags
```

## Why?

Why another language? Why another operating system?

All code dies. Some dies in hours, some dies in years. Seldom lives a
long life.

As a result, programmers keep writing essentially the same code again
and again, in different times, maybe under different stories.

This makes no sense. This is a waste of our lives. Code should live
forever.

Some code dies because they are no longer needed; that's fine. But
some code dies because it is too complex to maintain and modify;
that's bad. Often times, it is not because the algorithm is too
complex to comprehend, but the architecture losses its structure due
to unmanaged code debt.

Many programming langauges and systems put coding efficiency,
performance and safety as first considerations, but readability and
comprehensibility as the last 
([for example](http://andrewkelley.me/post/intro-to-zig.html)). 
I do not agree with those priorities. Good code quality fundamentally
comes from relentless iterations, and iterations require good code
understandings by *human* programmers. In the long run, code
comprehensibility dominates all.

E8VM proposes an entire set of langauge tool-chain that puts code
readability, or more precisely, code comprehensibility as the one and
only first priority.

To achieve high comprehensibility:

- It uses a simple programming language with only a small set of
  language features.
- It uses small files (80 chars per line, 300 lines maximum), and
  forbids circular dependencies among them, compiler enforced.
- It minimizes tool-chain setup frictions; it compiles and runs
  right inside the Web browsers.

These principles do not immediately achieve the best possible code
comprehensibility, but they at least encourages code reading and with
some proper execution, they will hopefully close the feedback loop on
code comprehensibility.

Long live our code.

## Copyright and License

The project developers own the copyright; my employer (Google) does
not own the copyright. Apache is the License.
