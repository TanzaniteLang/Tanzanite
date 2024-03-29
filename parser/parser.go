package parser

import (
    "strconv"
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/lexer"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/env"
)

type Token struct {
    Info tokens.Token
    Position tokens.Position
    Text string
}

type Parser struct {
    tokens []Token
    env env.Environment
}

func NewParser() *Parser {
    return &Parser{
        tokens: make([]Token, 0),
        env: env.NewEnv(),
    }
}

func (p *Parser) notEof() bool {
    return p.tokens[0].Info != tokens.Eof
}

func (p *Parser) current() Token {
    return p.tokens[0]
}

func (p *Parser) consume() Token {
    prev, tokens2 := p.tokens[0], p.tokens[1:]
    p.tokens = tokens2

    return prev
}

func (p *Parser) ProduceAST(code string) ast.Program {
    lex := lexer.InitLexer(code)

    for {
        pos, tok, text := lex.Lex()

        p.tokens = append(p.tokens, Token { Info: tok, Position: pos, Text: text})

        if tok == tokens.Eof {
            break
        }
    }

    prog := ast.Program {Body: make([]ast.Statement, 0)}

    for p.notEof() {
        prog.Body = append(prog.Body, p.parseStatement(&p.env))
    }

    return prog
}

func (p *Parser) parseStatement(e *env.Environment) ast.Statement {
    switch p.current().Info {
    case tokens.Def:
        panic("Def functions are not yet implemented!")
    case tokens.Fun:
        p.consume()
        fn := p.parseFunction(true).(ast.FunctionDecl)
        e.Fns[fn.Name] = &fn

        return fn
    case tokens.Return:
        p.consume()
        return ast.ReturnExpr{
            Value: p.parseExpression(),
        }
    case tokens.Identifier:
        fn, ok := e.Fns[p.current().Text]
        if ok { // This is a function call
            return p.parseFnCall(fn)
        }

        val, ok := e.Vars[p.current().Text]
        if !ok {
            stmt := p.parseVarDeclaration().(ast.VarDeclaration)
            e.Vars[stmt.Name] = &stmt
            return stmt
        }

        p.consume()
        p.consume()
        return ast.AssignExpr{
            Name: val,
            Value: p.parseExpression(),
        }
    default:
        return p.parseExpression()
    }
}

func (p *Parser) parseExpression() ast.Expression {
    return p.parseAdditiveExpr()
}

func (p *Parser) parseFnCall(fndecl *ast.FunctionDecl) ast.FunctionCall {
    calle := p.parseExpression()
    args := make([]ast.Expression, 0)

    current := p.current().Info
    for current != tokens.Semicolon {
        args = append(args, p.parseExpression())
        current = p.consume().Info
    }

    return ast.FunctionCall{
        Calle: calle,
        Args: args,
    }
}

func (p *Parser) parseFunction(isFun bool) ast.Statement {
    name := p.consume()
    current := p.consume()
    args := make([]ast.Statement, 0)
    returnType := make([]ast.Statement, 0)

    for current.Info != tokens.RBracket {
        if current.Info == tokens.Dot {
            p.consume()
            args = append(args, ast.VariadicArg{})
        } else {
            args = append(args, p.parseVarDeclaration())
        }
        current = p.consume()
    }

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
    }

    current = p.current()
    for current.Info != tokens.End {
        fn.Body = append(fn.Body, p.parseStatement(&p.env))
        current = p.current()
    }
    p.consume()

    return fn
}

func (p *Parser) parseVarDeclaration() ast.Statement {
    ident := p.consume()

    if p.current().Info == tokens.Colon {
        p.consume()

        varType := p.parseType()

        if p.current().Info == tokens.Assign {
            p.consume()

            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: p.parseExpression(),
            }
        } else {
            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType,
                Value: nil,
            }
        }
    } else if p.current().Info == tokens.Assign {
        p.consume()
        return ast.VarDeclaration{
            Name: ident.Text,
            Type: make([]ast.Statement, 0),
            Value: p.parseExpression(),
        }
    }

    return nil
}

func (p *Parser) parseType() []ast.Statement {
    typeConstruct := make([]ast.Statement, 0)
    typeConstruct = append(typeConstruct, ast.Identifier{Symbol: p.consume().Text})

    current := p.current().Info
    for current == tokens.Asterisk || current == tokens.DoubleAsterisk {
        if current == tokens.DoubleAsterisk {
            typeConstruct = append(typeConstruct, ast.Pointer{})
            typeConstruct = append(typeConstruct, ast.Pointer{})
        } else {
            typeConstruct = append(typeConstruct, ast.Pointer{})
        }
        p.consume()

        current = p.current().Info
    }

    return typeConstruct
}

func (p *Parser) parseMultiplicativeExpr() ast.Expression {
    left := p.parsePrimaryExpr()

    for p.current().Text == "/" || p.current().Text == "*" ||
        p.current().Text == "%" {
        operator := p.consume().Text
        right := p.parsePrimaryExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parseAdditiveExpr() ast.Expression {
    left := p.parseMultiplicativeExpr()

    for p.current().Text == "+" || p.current().Text == "-" {
        operator := p.consume().Text
        right := p.parseMultiplicativeExpr()
        left = ast.BinaryExpr{
            Left: left,
            Right: right,
            Operator: operator,
        }
    }

    return left
}

func (p *Parser) parsePrimaryExpr() ast.Expression {
    tok := p.current().Info;

    switch tok {
    case tokens.Identifier:
        return ast.Identifier{ Symbol: p.consume().Text }
    case tokens.String:
        return ast.String{ Value: p.consume().Text }
    case tokens.Float:
        val, _ := strconv.ParseFloat(p.consume().Text, 64)
        return ast.FloatLiteral{
            Value: val,
        }
    case tokens.Int:
        val, _ := strconv.ParseInt(p.consume().Text, 10, 64)
        return ast.IntLiteral{
            Value: val,
        }
    case tokens.LBracket:
        p.consume()
        defer p.consume()
        
        return p.parseExpression()
    }

    return nil
}
