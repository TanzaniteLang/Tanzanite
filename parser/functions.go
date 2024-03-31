package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

func (p *Parser) variadicCall(fndecl *ast.FunctionDecl) []ast.Expression {
    args := make([]ast.Expression, 0)
    argCount := len(fndecl.Arguments) - 1 // because variadic

    if p.current().Info != tokens.LBracket {
        panic("Missing (")
    }
    p.consume()

    for i := 0; i < argCount; i++ {
        args = append(args, p.parseExpression())

        if i + 1 == argCount {
            break 
        } else if p.current().Info == tokens.RBracket && i + 1 < argCount {
            panic("Invalid arg count")
        }

        p.consume()
    }

    for p.current().Info == tokens.Comma {
        p.consume()
        args = append(args, p.parseExpression())
    }

    if p.current().Info != tokens.RBracket {
        panic("Missing )")
    }

    p.consume()

    return args
}

func (p *Parser) functionCall(fndecl *ast.FunctionDecl) []ast.Expression {
    args := make([]ast.Expression, 0)
    argCount := len(fndecl.Arguments)

    if p.current().Info != tokens.LBracket {
        panic("Missing (")
    }
    p.consume()

    for i := 0; i < argCount; i++ {
        args = append(args, p.parseExpression())

        if i + 1 == argCount {
            break 
        } else if p.current().Info == tokens.RBracket && i + 1 < argCount {
            panic("Invalid arg count")
        }

        p.consume()
    }

    if p.current().Info != tokens.RBracket {
        panic("Missing )")
    }

    p.consume()

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
    p.consume()

    args := make([]ast.Statement, 0)
    returnType := make([]ast.Statement, 0)

    variadic := false

    current := p.current()
    for current.Info != tokens.RBracket {
        if current.Info == tokens.Identifier {
            args = append(args, p.parseVarDeclaration())
        } else if current.Info == tokens.Dot {
            p.consume()
            p.consume()
            p.consume()
            args = append(args, ast.VariadicArg{})
            variadic = true
        }
        current = p.current()

        if current.Info != tokens.Comma && current.Info != tokens.RBracket {
            panic("Missing ,")
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
    }

    fn := ast.FunctionDecl {
        Name: name.Text,
        Arguments: args,
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
