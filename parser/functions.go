package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) parseFnCall(fndecl *ast.FunctionDecl) ast.FunctionCall {
    p.parsingFn = true
    calle_pos := p.current().Position
    calle := p.parseExpression()
    p.parsingFn = false
    args := make([]ast.Expression, 0)

    needBracket := p.requireBrackets

    argCount := len(fndecl.Arguments)
    if fndecl.Variadic {
        argCount--
    }

    if needBracket && p.current().Info != tokens.LBracket {
        c := p.current()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Function call here is required to have (!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add (",
            Code: "(",
        })
        p.Dead = true
    }

    if p.current().Info == tokens.LBracket {
        needBracket = true
        p.consume()
    }

    if argCount > 0 {
        if !needBracket && calle_pos.Line != p.current().Position.Line {
            c := p.current()
            dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
            dbg.ThrowError("Function args must start on the same line as calle!", p.warn || p.Dead, nil)
            p.Dead = true
            return ast.FunctionCall{
                Calle: calle,
                Args: args,
            }
        }

        for {
            p.requireBrackets = true
            expr := p.parseExpression()
            p.requireBrackets = false
            if expr == nil {
                break
            }
            args = append(args, expr)

            if p.current().Info != tokens.Comma {
                break
            }

            p.consume()
        }
    }

    if needBracket && p.current().Info != tokens.RBracket {
        c := p.previous()
        dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column)
        dbg.ThrowError("Function call needs to close with )!", p.warn || p.Dead, &debug.Hint{
            Msg: "Add )",
            Code: ")",
        })
        p.Dead = true
    } else if needBracket && p.current().Info == tokens.RBracket {
        p.consume()
    }

    return ast.FunctionCall{
        Calle: calle,
        Args: args,
        Position: calle_pos,
    }
}

func (p *Parser) parseFunction(isFun bool) ast.Statement {
    start_pos := p.consume().Position
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
        Body: ast.Body{
            Scope: map[string]*ast.VarDeclaration{},
            Body: []ast.Statement{},
        },
        Variadic: variadic,
        Position: start_pos,
    }

    p.AppendScope(&fn.Body)

    current = p.current()
    for current.Info != tokens.End {
        fn.Body.Append(p.parseStatement())
        current = p.current()
    }
    p.consume()

    p.PopScope()

    return fn
}
