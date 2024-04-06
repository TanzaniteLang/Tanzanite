package analyzer

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "codeberg.org/Tanzanite/Tanzanite/ast"
)

func (a *Analyzer) analyzeAssignment(a2 *ast.AssignExpr) string {
    nameType := a.checkExpression(&a2.Name)
    valueType := a.checkExpression(&a2.Value)

    if nameType != valueType && !a.unaryStatment {
        dbg := debug.NewSourceLocation(a.Source, a2.Position.Line, a2.Position.Column)
        dbg.ThrowError(fmt.Sprintf("Type '%s' cannot be assigned to variable of type '%s'",
        valueType, nameType), a.Dead, nil)
        a.Dead = true
        return ""
    }

    return nameType
}
