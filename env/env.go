package env

import (
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

type Environment struct {
    // TODO: Once structs are things, string -> (whatever ast node)
    Vars map[string]*ast.VarDeclaration
    Fns map[string]*ast.FunctionDecl // TODO: Same for fns
}

func NewEnv() Environment {
    return Environment {
        Vars: make(map[string]*ast.VarDeclaration),
        Fns: make(map[string]*ast.FunctionDelc),
    }
}
