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

    pos int
    source string
    parsingFn bool // Breaks reccursion

    // Programmer info
    warn bool
    Dead bool
}

func NewParser(file string) *Parser {
    return &Parser{
        tokens: make([]Token, 0),
        env: env.NewEnv(),
        parsingFn: false,
        source: file,
        Dead: false,
        warn: false,
        pos: 0,
    }
}

func (p *Parser) notEof() bool {
    return p.tokens[p.pos].Info != tokens.Eof
}

func (p *Parser) current() Token {
    return p.tokens[p.pos]
}

func (p *Parser) consume() Token {
    p.pos++
    return p.tokens[p.pos - 1]
}

func (p *Parser) previous() Token {
    if p.pos == 0 {
        return p.tokens[0]
    }

    return p.tokens[p.pos - 1]
}

func (p *Parser) skipToNewLine() {
    pos := p.previous().Position.Line

    for p.current().Position.Line == pos {
        p.pos++
    }
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
        prog.Debug = append(prog.Debug, debug.NewSourceLocation(p.source, p.current().Position.Line, p.current().Position.Column))
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
    case tokens.Break:
        p.consume()
        return ast.LoopControlStatement{
            Break: true,
        }
    case tokens.Next:
        p.consume()
        return ast.LoopControlStatement{
            Break: false,
        }
    case tokens.If:
        return p.parseIf(false)
    case tokens.Unless:
        return p.parseIf(true)
    case tokens.While:
        return p.parseWhile(false)
    case tokens.Until:
        return p.parseWhile(true)
    case tokens.Return:
        p.consume()
        expr := p.parseExpression()
        if expr == nil {
            c := p.previous()
            dbg := debug.NewSourceLocation(p.source, c.Position.Line, c.Position.Column + 1 + uint64(len(c.Text)))
            dbg.ThrowError("Return statement is missing a value!", p.warn || p.Dead, nil)
            p.Dead = true
            p.skipToNewLine()
        }

        return ast.ReturnExpr{
            Value: expr,
        }
    case tokens.Identifier:
        fn, ok := p.env.Fns[p.current().Text]
        if ok { // This is a function call
            return p.parseFnCall(fn)
        }

        _, ok = p.env.Vars[p.current().Text]
        if !ok {
            possiblyVar := p.parseVarDeclaration()
            if possiblyVar == nil {
                return nil
            }
            stmt := possiblyVar.(ast.VarDeclaration)
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
