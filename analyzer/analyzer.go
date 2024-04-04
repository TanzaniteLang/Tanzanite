package analyzer

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

type Analyzer struct {
    Parser *parser.Parser
    Program *ast.Program
    Dead bool
    Source string
    Scopes []*ast.Body
}

func (a *Analyzer) findVariable(name string) *ast.VarDeclaration {
    last := len(a.Scopes) - 1
    for last >= 0 {
        decl, ok := a.Scopes[last].Scope[name]
        if ok {
            return decl
        }
        last--
    }
    return nil
}

func (a *Analyzer) AppendScope(scope *ast.Body) {
    a.Scopes = append(a.Scopes, scope)
}

func (a *Analyzer) PopScope() {
    a.Scopes = a.Scopes[:len(a.Scopes) - 1]
}

func (a *Analyzer) checkVariable(v *ast.VarDeclaration) {
    if len(v.Type) > 0 {
        resultType := a.checkExpression(&v.Value)
        if resultType != ast.StrType(v.Type) {
            dbg := debug.NewSourceLocation(a.Source, v.Position.Line, v.Position.Column)
            dbg.ThrowError(fmt.Sprintf("Type '%s' cannot be assigned to variable of type '%s'",
            resultType, ast.StrType(v.Type)), a.Dead, nil)
            a.Dead = true
        }
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

        decl := a.findVariable(i.Symbol)
        if decl == nil {
            dbg := debug.NewSourceLocation(a.Source, i.Position.Line, i.Position.Column)
            dbg.ThrowError("Unknown variable '" + i.Symbol + "'!", a.Dead, nil)
            a.Dead = true
            return ""
        }

        return ast.StrType(decl.Type)
    case ast.BinaryExprType:
        b := (*expr).(ast.BinaryExpr)
        left_expr := a.checkExpression(&b.Left)
        right_expr := a.checkExpression(&b.Right)
        if left_expr != right_expr {
            dbg := debug.NewSourceLocation(a.Source, b.Position.Line, b.Position.Column)
            dbg.ThrowError(fmt.Sprintf("Cannot perform '%s' on types '%s' and '%s'",
            b.Operator, left_expr, right_expr), a.Dead, nil)
            a.Dead = true
        }

        return left_expr 
    }

    return ""
}

func (a *Analyzer) Analyze() {
    a.AppendScope(&a.Program.Body)
    for _, stmt := range a.Program.Body.Body {
        switch stmt.GetKind() {
        case ast.VarDeclarationType:
            v := stmt.(ast.VarDeclaration)
            a.checkVariable(&v)
        case ast.FunctionDeclType:
            _ = stmt.(ast.FunctionDecl)
            
        }
    }
}
