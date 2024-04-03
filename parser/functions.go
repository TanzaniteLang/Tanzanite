package parser

import (
    "fmt"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) variadicCall(fndecl *ast.FunctionDecl) []ast.Expression {
    args := make([]ast.Expression, 0)
    needBracket := p.requireBrackets
    argCount := len(fndecl.Arguments) - 1 // because variadic

    if needBracket && p.current().Info == tokens.LBracket {
        p.consume()
    } else if needBracket && p.current().Info != tokens.LBracket {
        c := p.current()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Function call here is required to have (!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add (",
            Code: "(",
        })
        p.Dead = true
    } else if p.current().Info == tokens.LBracket {
        needBracket = true
        p.consume()
    }

    for i := 0; i < argCount; i++ {
        p.requireBrackets = true
        expr := p.parseExpression()
        if expr == nil {
            debug.LogError("Invalid argument count for function \"" + fndecl.Name + "\"!", &debug.Hint{
                Msg: fmt.Sprintf("Function requires %d arguments", argCount),
            })
            p.Dead = true
        }
        args = append(args, expr)
        p.requireBrackets = false

        if i + 1 == argCount {
            break 
        }

        p.consume()
    }

    for p.current().Info == tokens.Comma {
        p.consume()
        p.requireBrackets = true
        args = append(args, p.parseExpression())
        p.requireBrackets = false
    }

    if needBracket && p.current().Info == tokens.RBracket {
        p.consume()
    } else if needBracket && p.current().Info != tokens.RBracket {
        c := p.previous()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1)
        dbg.ThrowError("Function call is missing )!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add )",
            Code: ")",
        })
        p.Dead = true
        p.skipToNewLine()
    }

    return args
}

func (p *Parser) functionCall(fndecl *ast.FunctionDecl) []ast.Expression {
    args := make([]ast.Expression, 0)
    needBracket := p.requireBrackets
    argCount := len(fndecl.Arguments)

    if needBracket && p.current().Info == tokens.LBracket {
        p.consume()
    } else if needBracket && p.current().Info != tokens.LBracket {
        c := p.current()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Function call here is required to have (!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add (",
            Code: "(",
        })
        p.Dead = true
    } else if p.current().Info == tokens.LBracket {
        needBracket = true
        p.consume()
    }

    for i := 0; i < argCount; i++ {
        p.requireBrackets = true
        expr := p.parseExpression()
        e := fndecl.Arguments[i].(ast.VarDeclaration)

        if expr == nil && e.Value == nil {
            debug.LogError("Invalid argument count for function \"" + fndecl.Name + "\"!", &debug.Hint{
                Msg: fmt.Sprintf("Function requires %d arguments", argCount),
            })
            p.Dead = true
        }

        if expr == nil && e.Value != nil {
            args = append(args, e.Value)
            p.requireBrackets = false
            if i + 1 == argCount {
                break 
            }
            continue
        } else {
            args = append(args, expr)
        }
        p.requireBrackets = false

        if i + 1 == argCount {
            break 
        }

        if p.current().Info == tokens.Comma {
            p.consume()
        }
    }

    if needBracket && p.current().Info == tokens.RBracket {
        p.consume()
    } else if needBracket && p.current().Info != tokens.RBracket {
        c := p.previous()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1)
        dbg.ThrowError("Function call is missing )!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add )",
            Code: ")",
        })
        p.Dead = true
        p.skipToNewLine()
    }

    return args
}

func (p *Parser) parseFnCall(fndecl *ast.FunctionDecl) ast.FunctionCall {
    p.parsingFn = true
    calle := p.parseExpression()
    p.parsingFn = false
    args := make([]ast.Expression, 0)

    if fndecl.Variadic {
        args = p.variadicCall(fndecl)
    } else {
        args = p.functionCall(fndecl)
    }

    return ast.FunctionCall{
        Calle: calle,
        Args: args,
    }
}

func (p *Parser) parseFunction(isFun bool) ast.Statement {
    name := p.consume()
    if p.current().Info != tokens.LBracket {
        c := p.current()

        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Function arguments must be in ()!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add missing (", 
            Code: "(",
        })
        p.Dead = true
    } else {
        p.consume()
    }

    fail := false

    args := make([]ast.Statement, 0)
    returnType := make([]ast.Statement, 0)

    variadic := false

    current := p.current()
    for current.Info != tokens.RBracket {
        if current.Info == tokens.Identifier {
            args = append(args, p.parseVarDeclaration())
        } else if current.Info == tokens.Dot {
            p.consume()
            if p.current().Info != tokens.Dot {
                c := p.current()

                dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
                dbg.ThrowError("Variadic arg needs 3 dots, got 1!", p.warn || p.Dead, nil)
                p.Dead = true
                fail = true
                break
            }
            p.consume()
            if p.current().Info != tokens.Dot {
                c := p.current()

                dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
                dbg.ThrowError("Variadic arg needs 3 dots, got 2!", p.warn || p.Dead, nil)
                p.Dead = true
                fail = true
                break
            }
            p.consume()
            args = append(args, ast.VariadicArg{})
            variadic = true
        }
        current = p.current()

        if current.Info != tokens.Comma && current.Info != tokens.RBracket {
            c := p.current()

            dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
            dbg.ThrowError("Expected , or ) but got " + current.Text + " instead!", p.warn || p.Dead, nil)
            p.Dead = true
            p.skipToNewLine()
            break
        } else if current.Info == tokens.RBracket {
            break
        }
        p.consume()

        current = p.current()
    }
    p.consume()

    if p.current().Info == tokens.Colon {
        p.consume()
        returnType = p.parseType()
    } else {
        c := p.previous()

        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1)
        dbg.ThrowWarning("No explicit return type specified!", p.warn || p.Dead, &debug.Hint{
            Msg: "Void will be used as return type, if you don't want that, provide a type", 
            Code: ": Type",
        })
        p.warn = true
        p.skipToNewLine()
        returnType = append(returnType, ast.TypeLiteral{
            Type: []ast.Statement{ast.Identifier{
                Symbol: "void",
            }},
        })
    }

    fn := ast.FunctionDecl {
        Name: name.Text,
        Arguments: args,
        Failed: fail,
        ReturnType: returnType,
        Immutable: isFun,
        Body: make([]ast.Statement, 0),
        Variadic: variadic,
    }

    current = p.current()
    for current.Info != tokens.End {
        fn.Debug = append(fn.Debug, debug.NewSourceLocation(p.source, current.Position.Line, current.Position.Column))
        fn.Body = append(fn.Body, p.parseStatement())
        current = p.current()
    }
    p.consume()

    return fn
}
