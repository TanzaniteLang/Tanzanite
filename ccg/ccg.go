package ccg

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

type Source struct {
    Name string
    Functions []ast.FunctionDecl

    source string
}

func NewSource(name string) *Source {
    return &Source{
        Name: name,
        Functions: make([]ast.FunctionDecl, 0),
        source: "",
    }
}

func (s *Source) Generate() string {
    for _, fn := range s.Functions {
        s.source += fn.GenDecl() + ";\n"
    }

    s.source += "\n"

    for _, fn := range s.Functions {
        if len(fn.Body) > 0 {
            s.source += fn.GenDecl() + " {\n" + fn.GenBody() + "}\n"
        }
    }

    return s.source
}
