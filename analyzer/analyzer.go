package analyzer

import (
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "github.com/gookit/goutil/dump"
)

type Analyzer struct {
    Parser *parser.Parser
    Program *ast.Program
}

func (a *Analyzer) checkVariable(v *ast.VarDeclaration) {
    if len(v.Type) > 0 {
        // TODO: Stringify the type
        dump.Println(v.Value)
        dump.Println(a.checkExpression(&v.Value))
    }
}

func (a *Analyzer) checkExpression(expr *ast.Expression) string {
    switch (*expr).GetKind() {
    case ast.IntLiteralType:
        return "int"
    case ast.FloatLiteralType:
        return "float"
    case ast.CharType:
        return "char"
    case ast.BoolType:
        return "bool"
    case ast.IdentifierType:
        i := (*expr).(ast.Identifier)
        return i.Symbol
    case ast.BinaryExprType:
        b := (*expr).(ast.BinaryExpr)
        if a.checkExpression(&b.Left) != a.checkExpression(&b.Right) {
            panic("Types do not match!")
        }

        return a.checkExpression(&b.Left)
    }

    return ""
}

func (a *Analyzer) Analyze() {
    for _, stmt := range a.Program.Body.Body {
        switch stmt.GetKind() {
        case ast.VarDeclarationType:
            v := stmt.(ast.VarDeclaration)
            a.checkVariable(&v)
        }
    }
}
