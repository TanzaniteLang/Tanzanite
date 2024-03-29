package main

import (
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "github.com/gookit/goutil/dump"
)

func main() {
    par := parser.NewParser()
    out := par.ProduceAST(`fun printf(format: Char*, ...): Int
end

fun main(argc: Int, argv: Char**): Int
    printf "Hello World! Argc: %d, Argv: %p\n", argc, argv; # Variadic function
    return 10 - 10
end
`)

    dump.Config(func (o *dump.Options) {
        o.MaxDepth = 10
    })

    dump.Println(out)
}
