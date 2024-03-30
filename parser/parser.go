package parser

import (
    "codeberg.org/Tanzanite/Tanzanite/tokens"
    "codeberg.org/Tanzanite/Tanzanite/lexer"
    "codeberg.org/Tanzanite/Tanzanite/ast"
    "codeberg.org/Tanzanite/Tanzanite/env"
    "codeberg.org/Tanzanite/Tanzanite/debug"
)

type Token struct {
    Info tokens.Token
    Position tokens.Position
    Text string
}

type Parser struct {
    tokens []Token
    env env.Environment

    source string
    parsingFn bool // Breaks reccursion
}

func NewParser(file string) *Parser {
    return &Parser{
        tokens: make([]Token, 0),
        env: env.NewEnv(),
        parsingFn: false,
        source: file,
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
        prog.Debug = append(prog.Debug, debug.NewSourceLocation(p.source, p.current().Position.Line))
        prog.Body = append(prog.Body, p.parseStatement())
    }

    return prog
}

func (p *Parser) parseStatement() ast.Statement {
    switch p.current().Info {
    case tokens.Def:
        panic("Def functions are not yet implemented!")
    case tokens.Fun:
        p.consume()
        fn := p.parseFunction(true).(ast.FunctionDecl)
        p.env.Fns[fn.Name] = &fn

        return fn
    case tokens.Return:
        p.consume()
        return ast.ReturnExpr{
            Value: p.parseExpression(),
        }
    case tokens.Identifier:
        fn, ok := p.env.Fns[p.current().Text]
        if ok { // This is a function call
            return p.parseFnCall(fn)
        }

        _, ok = p.env.Vars[p.current().Text]
        if !ok {
            stmt := p.parseVarDeclaration().(ast.VarDeclaration)
            p.env.Vars[stmt.Name] = &stmt
            return stmt
        }
        return p.parseAssignExpr()
    default:
        return p.parseExpression()
    }
}

func (p *Parser) parseExpression() ast.Expression {
    return p.parseAssignExpr()
}
