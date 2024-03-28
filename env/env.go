package env

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

type Environment struct {
    Vars map[string]*ast.VarDeclaration
}

func NewEnv() Environment {
    return Environment {
        Vars: make(map[string]*ast.VarDeclaration),
    }
}
