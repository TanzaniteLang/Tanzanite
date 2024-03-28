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
    case tokens.Identifier:
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

func (p *Parser) parseVarDeclaration() ast.Statement {
    ident := p.consume()


    if p.current().Info == tokens.Colon {
        p.consume()

        varType := p.consume()

        if p.current().Info == tokens.Assign {
            p.consume()

            return ast.VarDeclaration{
                Name: ident.Text,
                Type: varType.Text,
                Value: p.parseExpression(),
            }
        }
    } else if p.current().Info == tokens.Assign {
        p.consume()
        return ast.VarDeclaration{
            Name: ident.Text,
            Type: "??",
            Value: p.parseExpression(),
        }
    }

    return nil
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
