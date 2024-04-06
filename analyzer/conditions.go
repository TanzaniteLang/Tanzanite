package analyzer

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

func (a *Analyzer) analyzeIf(i *ast.IfStatement) {
    resultType := a.checkExpression(&i.Condition)

    if resultType != "bool" {
        dbg := debug.NewSourceLocation(a.Source, i.Position.Line, i.Position.Column)
        dbg.ThrowError(fmt.Sprintf("Expected type 'bool', got '%s'!", resultType), a.Dead, nil)
        a.Dead = true
    }

    for _, stmt := range i.Body.Body {
        a.analyzeStatement(&stmt)
    }

    if i.Next != nil {
        if i.Next.GetKind() == ast.ElsifStatementType {
            e := i.Next.(ast.ElsifStatement)
            a.analyzeElsif(&e)
        } else {
            e := i.Next.(ast.ElseStatement)

            for _, stmt := range e.Body.Body {
                a.analyzeStatement(&stmt)
            }
        }
    }
}

func (a *Analyzer) analyzeElsif(e *ast.ElsifStatement) {
    resultType := a.checkExpression(&e.Condition)

    if resultType != "bool" {
        dbg := debug.NewSourceLocation(a.Source, e.Position.Line, e.Position.Column)
        dbg.ThrowError(fmt.Sprintf("Expected type 'bool', got '%s'!", resultType), a.Dead, nil)
        a.Dead = true
    }

    for _, stmt := range e.Body.Body {
        a.analyzeStatement(&stmt)
    }

    if e.Next != nil {
        if e.Next.GetKind() == ast.ElsifStatementType {
            e2 := e.Next.(ast.ElsifStatement)
            a.analyzeElsif(&e2)
        } else {
            e2 := e.Next.(ast.ElseStatement)

            for _, stmt := range e2.Body.Body {
                a.analyzeStatement(&stmt)
            }
        }
    }
}
