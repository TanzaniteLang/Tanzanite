package analyzer

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (a *Analyzer) analyzeWhile(w *ast.WhileStatement) {
    resultType := a.checkExpression(&w.Condition)

    if resultType != "bool" {
        dbg := debug.NewSourceLocation(a.Source, w.Position.Line, w.Position.Column)
        dbg.ThrowError(fmt.Sprintf("Expected type 'bool', got '%s'!", resultType), a.Dead, nil)
        a.Dead = true
    }

    for _, stmt := range w.Body.Body {
        a.analyzeStatement(&stmt)
    }
}
