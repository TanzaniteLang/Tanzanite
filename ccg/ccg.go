package ccg

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

import _ "embed"

//go:embed __bootstrap.c
var bootstrap string

type Source struct {
    Name string
    Functions []ast.FunctionDecl

    source string
}

func NewSource(name string) *Source {
    return &Source{
        Name: name,
        Functions: make([]ast.FunctionDecl, 0),
        // Tanzanite boilerplate
        source: bootstrap,
    }
}

func (s *Source) Generate() string {
    for _, fn := range s.Functions {
        s.source += fn.StringifyHead() + ";\n"
    }

    s.source += "\n"

    for _, fn := range s.Functions {
        if len(fn.Body) > 0 {
            s.source += fn.Stringify() + "\n"
        }
    }

    return s.source[:len(s.source) - 2]
}
