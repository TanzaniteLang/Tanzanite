package analyzer

import (
    "fmt"
    "strings"
    "codeberg.org/Tanzanite/Tanzanite/parser"
    "codeberg.org/Tanzanite/Tanzanite/debug"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "github.com/gookit/goutil/dump"
)

type Analyzer struct {
    Parser *parser.Parser
    Program *ast.Program
    Dead bool
    Source string
    Scopes []*ast.Body

    unaryStatment bool
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
    if len(v.Type) > 0 && v.Value != nil {
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
    case ast.StringType:
        return "char*"
    case ast.BoolType:
        return "bool"
    case ast.UnaryExprType:
        u := (*expr).(ast.UnaryExpr)
        resultType := a.checkExpression(&u.Operand)

        switch u.Operator {
        case "*":
            if strings.HasSuffix(resultType, "*") {
                return resultType[:len(resultType) - 1]
            } else {
                dbg := debug.NewSourceLocation(a.Source, u.Position.Line, u.Position.Column)
                dbg.ThrowError(fmt.Sprintf("Cannot dereference type '%s'!", resultType), a.Dead, nil)
                a.Dead = true
            }
        case "&":
            return resultType + "*"
        case "!":
            if resultType != "bool" {
                dbg := debug.NewSourceLocation(a.Source, u.Position.Line, u.Position.Column)
                dbg.ThrowError(fmt.Sprintf("Expected type 'bool', got '%s'!", resultType), a.Dead, nil)
                a.Dead = true
                return ""
            }
            return "bool"
        }

        return resultType
    case ast.IdentifierType:
        i := (*expr).(ast.Identifier)

        if i.Symbol == "nil" {
            return "void*"
        }

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

        switch b.Operator {
        case "==":
            return "bool"
        case "!=":
            return "bool"
        case "<":
            return "bool"
        case "<=":
            return "bool"
        case ">":
            return "bool"
        case ">=":
            return "bool"
        case "||":
            return "bool"
        case "&&":
            return "bool"
        }

        return left_expr 
    case ast.FunctionCallType:
        f := (*expr).(ast.FunctionCall)
        return a.analyzeFnCall(&f)
    case ast.TypeCastType:
        t := (*expr).(ast.TypeCast)

        return ast.StrType(t.Target.Type)
    case ast.BracketExprType:
        b := (*expr).(ast.BracketExpr)

        return a.checkExpression(&b.Expr)
    case ast.AssignExprType:
        a2 := (*expr).(ast.AssignExpr)
        return a.analyzeAssignment(&a2)
    }

    return ""
}

func (a *Analyzer) analyzeFn(fndecl *ast.FunctionDecl) {
    a.AppendScope(&fndecl.Body)

    need_default_value := false

    for _, arg := range fndecl.Arguments {
        if arg.GetKind() == ast.VarDeclarationType {
            v := arg.(ast.VarDeclaration)
            if need_default_value && v.Value == nil {
                dbg := debug.NewSourceLocation(a.Source, v.Position.Line, v.Position.Column)
                dbg.ThrowError("Optional function parameter '" + v.Name + "' is expected to have a value!", a.Dead, nil)
                fndecl.Failed = true
                a.Dead = true
                break
            }

            if v.Value != nil {
                need_default_value = true
                a.checkVariable(&v)
            }
        }
    }

    for  _, stmt := range fndecl.Body.Body {
        if stmt.GetKind() == ast.ReturnExprType {
            ret := stmt.(ast.ReturnExpr)
            resultType := a.checkExpression(&ret.Value)
            if resultType != ast.StrType(fndecl.ReturnType) {
                dbg := debug.NewSourceLocation(a.Source, ret.Position.Line, ret.Position.Column + uint64(len("return ")))
                dbg.ThrowError(fmt.Sprintf("Function returns type '%s', got '%s'",
                ast.StrType(fndecl.ReturnType), resultType), a.Dead, nil)
                a.Dead = true
            }
        } else {
            a.analyzeStatement(&stmt)
        }
    }

    a.PopScope()
}


func (a *Analyzer) analyzeFnCall(fncall *ast.FunctionCall) string {
    if a.Parser.Globals.HasFunction(fncall.Calle) {
        fn, _ := a.Parser.Globals.Scope[fncall.Calle]

        minArgs := 0
        maxArgs := 0

        for _, arg := range fn.Arguments {
            if arg.GetKind() == ast.VarDeclarationType {
                v := arg.(ast.VarDeclaration)
                if v.Value == nil {
                    minArgs++
                }
                maxArgs++
            }
        }

        if len(fncall.Args) < minArgs {
            dbg := debug.NewSourceLocation(a.Source, fncall.Position.Line, fncall.Position.Column)
            dbg.ThrowError(fmt.Sprintf("Function '%s' expects %d arguments, got %d!",
                           fncall.Calle, minArgs, len(fncall.Args)), a.Dead, nil)
            a.Dead = true
        }

        if len(fncall.Args) > minArgs {
            minArgs += len(fncall.Args)
        }

        for i := 0; i < minArgs; i++ {
            if i >= len(fn.Arguments) {
                continue
            }

            callarg := fncall.Args[i]
            resultType := a.checkExpression(&callarg)

            ar := fn.Arguments[i]
            if ar.GetKind() == ast.VariadicArgType {
                continue
            }

            arg := ar.(ast.VarDeclaration)

            if resultType != ast.StrType(arg.Type) {
                dbg := debug.NewSourceLocation(a.Source, arg.Position.Line, arg.Position.Column)
                dbg.ThrowError(fmt.Sprintf("%d. argument of the function '%s' expected type '%s' but got '%s' instead!",
                i + 1, fncall.Calle, ast.StrType(arg.Type), resultType), a.Dead, nil)
                a.Dead = true
            }
        }

        if minArgs < maxArgs {
            for i := minArgs; i < maxArgs; i++ {
                v := fn.Arguments[i].(ast.VarDeclaration)
                fncall.Args = append(fncall.Args, v.Value)
            }

        }

        return ast.StrType(fn.ReturnType)
    } else {
        dbg := debug.NewSourceLocation(a.Source, fncall.Position.Line, fncall.Position.Column)
        dbg.ThrowError("Unknown function '" + fncall.Calle + "'!", a.Dead, nil)
        a.Dead = true
    }

    return ""
}

func (a *Analyzer) analyzeStatement(stmt *ast.Statement) {
    switch (*stmt).GetKind() {
    case ast.VarDeclarationType:
        v := (*stmt).(ast.VarDeclaration)
        a.checkVariable(&v)
    case ast.FunctionDeclType:
        fn := (*stmt).(ast.FunctionDecl)
        a.analyzeFn(&fn)
    case ast.FunctionCallType:
        fn := (*stmt).(ast.FunctionCall)
        a.analyzeFnCall(&fn)
    case ast.WhileStatementType:
        w := (*stmt).(ast.WhileStatement)
        a.analyzeWhile(&w)
    case ast.LoopControlStatementType:
        return
    case ast.IfStatementType:
        i := (*stmt).(ast.IfStatement)
        a.analyzeIf(&i)
    case ast.AssignExprType:
        a2 := (*stmt).(ast.AssignExpr)
        a.analyzeAssignment(&a2)
    case ast.UnaryExprType:
        // XXX: This ain't good, but ideal for now!
        u := (*stmt).(ast.UnaryExpr)
        e := (*stmt).(ast.Expression)
        a.unaryStatment = u.Operator == "*"
        a.checkExpression(&e)
        a.unaryStatment = false
    default:
        dump.Println(stmt)
    }
}

func (a *Analyzer) Analyze() {
    a.AppendScope(&a.Program.Body)
    for _, stmt := range a.Program.Body.Body {
        a.analyzeStatement(&stmt)
    }
}
