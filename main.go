package main

import (
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "github.com/gookit/goutil/dump"
)

func main() {
    par := parser.NewParser()
    out := par.ProduceAST(`ahoj = (7 + 4) * 2
ahoj = 4`)

    dump.Config(func (o *dump.Options) {
        o.MaxDepth = 10
    })

    dump.Println(out)
}
